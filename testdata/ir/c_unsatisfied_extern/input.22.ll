; ModuleID = 'testdata/ir/c_unsatisfied_extern/source.c'
source_filename = "testdata/ir/c_unsatisfied_extern/source.c"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@_ZTVN10__cxxabiv117__class_type_infoE = external unnamed_addr constant { [0 x ptr] }
@_ZTVN10__cxxabiv120__si_class_type_infoE = external unnamed_addr constant { [0 x ptr] }
@str = private unnamed_addr constant [2 x i8] c"x\00"
@_ZTI = constant { ptr, ptr } { ptr getelementptr inbounds (ptr, ptr @_ZTVN10__cxxabiv117__class_type_infoE, i64 2), ptr @str }

declare void @__cxa_pure_virtual()

@vt = constant [2 x ptr] [ptr @__cxa_pure_virtual, ptr @_ZTI]

define i32 @main() {
entry:
  ret i32 0
}
