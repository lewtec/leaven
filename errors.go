package main

import "errors"

// Sentinel errors for non-wrapping failures (janitor: table then %w).
var (
	errCsmithConfig = errors.New("csmith configuration")

	errNonPointerSource     = errors.New("non-pointer source parameter")
	errMismatchedSrcElem    = errors.New("mismatched source and element types")
	errNonConstantStructIdx = errors.New("non-constant index into struct")
	errUnsupportedIndexType = errors.New("unsupported type to index into")

	errIncompatiblePointers   = errors.New("incompatible pointer types")
	errUnsupportedICmpPred    = errors.New("unsupported comparison predicate")
	errUnsupportedZextTo      = errors.New("unsupported To type for zext")
	errUnsupportedZextFrom    = errors.New("unsupported From type for zext")
	errMismatchedZextTypes    = errors.New("mismatched types for zext")
	errUnsupportedInstruction = errors.New("unsupported instruction type")
	errStructIndexRange       = errors.New("struct index out of range")
	errUnsupportedAggregate   = errors.New("unsupported aggregate type for index")
	errIntToPtr               = errors.New("converting an integer to a pointer violates Go's unsafe.Pointer rules")

	errUnsupportedTerminator = errors.New("unsupported block terminator type")

	errUnsupportedFloatType = errors.New("unsupported floating-point type")
	errUnsupportedType      = errors.New("unsupported type")

	errIntConstTooLarge     = errors.New("integer constant too large")
	errUnsupportedUndefType = errors.New("unsupported type for undefined constant")
	errUnsupportedValueType = errors.New("unsupported type of value to translate")
)
