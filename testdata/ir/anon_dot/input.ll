; clang-style nested anonymous struct (e.g. tree-sitter headers)
; ModuleID = 'anon_dot'
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.anon.1 = type { i8, i8 }

@g = internal constant %struct.anon.1 { i8 1, i8 2 }, align 1

define dso_local i8 @get() {
entry:
  %0 = load i8, i8* getelementptr inbounds (%struct.anon.1, %struct.anon.1* @g, i32 0, i32 0), align 1
  ret i8 %0
}
