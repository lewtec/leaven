; void mid-function return (haskell scanner pattern)
define void @f(i32 %x) {
entry:
  %cmp = icmp eq i32 %x, 0
  br i1 %cmp, label %early, label %rest

early:
  ret void

rest:
  %add = add i32 %x, 1
  ret void
}
