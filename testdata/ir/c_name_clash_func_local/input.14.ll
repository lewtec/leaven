target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

%struct.Interpolation = type { i32 }

define internal i32 @is_verbatim(%struct.Interpolation* %i) {
entry:
  %0 = getelementptr inbounds %struct.Interpolation, %struct.Interpolation* %i, i32 0, i32 0
  %1 = load i32, i32* %0
  %and = and i32 %1, 1
  ret i32 %and
}

define dso_local i32 @scan(%struct.Interpolation* %cur) {
entry:
  %is_verbatim = alloca i32, align 4
  store i32 0, i32* %is_verbatim
  %0 = getelementptr inbounds %struct.Interpolation, %struct.Interpolation* %cur, i32 0, i32 0
  %1 = load i32, i32* %0
  %and = and i32 %1, 1
  %tobool = icmp ne i32 %and, 0
  br i1 %tobool, label %if.then, label %if.end

if.then:
  store i32 1, i32* %is_verbatim
  br label %if.end

if.end:
  %call = call i32 @is_verbatim(%struct.Interpolation* %cur)
  %tobool2 = icmp ne i32 %call, 0
  br i1 %tobool2, label %if.then3, label %if.end4

if.then3:
  ret i32 1

if.end4:
  %2 = load i32, i32* %is_verbatim
  ret i32 %2
}

define dso_local i32 @main() {
  %i = alloca %struct.Interpolation
  %f = getelementptr inbounds %struct.Interpolation, %struct.Interpolation* %i, i32 0, i32 0
  store i32 1, i32* %f
  %r = call i32 @scan(%struct.Interpolation* %i)
  %cmp = icmp ne i32 %r, 0
  %conv = zext i1 %cmp to i32
  %sub = sub i32 1, %conv
  ret i32 %sub
}
