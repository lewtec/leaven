; ModuleID = 'testdata/ir/cpp_dynamic_cast/source.cpp'
source_filename = "testdata/ir/cpp_dynamic_cast/source.cpp"
target datalayout = "e-m:e-p270:32:32-p271:32:32-p272:64:64-i64:64-i128:128-f80:128-n8:16:32:64-S128"
target triple = "x86_64-unknown-linux-gnu"

@_ZTVN10__cxxabiv117__class_type_infoE = external global [0 x ptr]
@_ZTVN10__cxxabiv120__si_class_type_infoE = external global [0 x ptr]
@_ZTS8ProbElem = constant [10 x i8] c"8ProbElem\00"
@_ZTI8ProbElem = constant { ptr, ptr } { ptr getelementptr inbounds (ptr, ptr @_ZTVN10__cxxabiv117__class_type_infoE, i64 2), ptr @_ZTS8ProbElem }
@_ZTS13GroupProbElem = constant [16 x i8] c"13GroupProbElem\00"
@_ZTI13GroupProbElem = constant { ptr, ptr, ptr } { ptr getelementptr inbounds (ptr, ptr @_ZTVN10__cxxabiv120__si_class_type_infoE, i64 2), ptr @_ZTS13GroupProbElem, ptr @_ZTI8ProbElem }
@_ZTV13GroupProbElem = constant [3 x ptr] [ptr null, ptr @_ZTI13GroupProbElem, ptr @noop]
@_ZTV8ProbElem = constant [3 x ptr] [ptr null, ptr @_ZTI8ProbElem, ptr @noop]
@ok = private unnamed_addr constant [3 x i8] c"ok\00"

define void @noop() {
  ret void
}

declare ptr @__dynamic_cast(ptr, ptr, ptr, i64)
declare i32 @puts(ptr)

define i32 @main() {
entry:
  %d = alloca [8 x i8], align 8
  store ptr getelementptr inbounds (ptr, ptr @_ZTV13GroupProbElem, i64 2), ptr %d, align 8
  %hit = call ptr @__dynamic_cast(ptr %d, ptr @_ZTI8ProbElem, ptr @_ZTI13GroupProbElem, i64 0)
  %miss_obj = alloca [8 x i8], align 8
  store ptr getelementptr inbounds (ptr, ptr @_ZTV8ProbElem, i64 2), ptr %miss_obj, align 8
  %miss = call ptr @__dynamic_cast(ptr %miss_obj, ptr @_ZTI8ProbElem, ptr @_ZTI13GroupProbElem, i64 0)
  %ok1 = icmp ne ptr %hit, null
  %ok2 = icmp eq ptr %miss, null
  %both = and i1 %ok1, %ok2
  br i1 %both, label %good, label %bad

good:
  %p = call i32 @puts(ptr @ok)
  ret i32 0

bad:
  ret i32 1
}
