package main

import "github.com/llir/llvm/ir"

func fixMalloc(f *ir.Func) {
	var idx Index
	idx.Add(f)

	for _, b := range f.Blocks {
		for _, inst := range b.Insts {
			switch inst := inst.(type) {
			case *ir.InstCall:
				callee, err := FormatValue(inst.Callee)
				if err != nil {
					continue
				}
				switch callee {
				case "malloc", "calloc":
					users := idx.Users(inst)
					if len(users) == 1 {
						// If the return value of malloc was immediately cast to another type,
						// tell our Malloc function to allocate that type instead of bytes.
						if bc, ok := users[0].(*ir.InstBitCast); ok {
							inst.Typ = bc.To
							idx.ReplaceValue(bc, inst)
							idx.DeleteInstruction(bc)
						}
					}

				case "free":
					// Leave free in the IR so codegen emits libc.Free (GC no-op).
					// Still drop a sole bitcast to *byte on the argument.
					if len(inst.Args) == 1 {
						if bc, ok := inst.Args[0].(*ir.InstBitCast); ok && len(idx.Users(bc)) == 1 {
							idx.ReplaceValue(bc, bc.From)
							idx.DeleteInstruction(bc)
						}
					}
				}
			}
		}
	}
}
