package cast_test

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/aura-studio/cast"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// =============================================================================
// Test Generators
// =============================================================================

// genInt8InRange generates int8 values within a safe range for round-trip testing
func genInt8InRange() gopter.Gen {
	return gen.Int8Range(math.MinInt8, math.MaxInt8)
}

// genInt16InRange generates int16 values within a safe range
func genInt16InRange() gopter.Gen {
	return gen.Int16Range(math.MinInt16, math.MaxInt16)
}

// genInt32InRange generates int32 values within a safe range
func genInt32InRange() gopter.Gen {
	return gen.Int32Range(math.MinInt32, math.MaxInt32)
}

// genInt64InRange generates int64 values
func genInt64InRange() gopter.Gen {
	return gen.Int64()
}

// genUint8 generates uint8 values
func genUint8() gopter.Gen {
	return gen.UInt8()
}

// genUint16 generates uint16 values
func genUint16() gopter.Gen {
	return gen.UInt16()
}

// genUint32 generates uint32 values
func genUint32() gopter.Gen {
	return gen.UInt32()
}

// genUint64 generates uint64 values
func genUint64() gopter.Gen {
	return gen.UInt64()
}

// genFloat32 generates float32 values (excluding NaN and Inf)
func genFloat32() gopter.Gen {
	return gen.Float32().SuchThat(func(f float32) bool {
		return !math.IsNaN(float64(f)) && !math.IsInf(float64(f), 0)
	})
}

// genFloat64 generates float64 values (excluding NaN and Inf)
func genFloat64() gopter.Gen {
	return gen.Float64().SuchThat(func(f float64) bool {
		return !math.IsNaN(f) && !math.IsInf(f, 0)
	})
}

// genNegativeInt generates negative int values
func genNegativeInt() gopter.Gen {
	return gen.IntRange(math.MinInt, -1)
}

// genNegativeFloat64 generates negative float64 values
func genNegativeFloat64() gopter.Gen {
	return gen.Float64Range(math.SmallestNonzeroFloat64*-1e100, -math.SmallestNonzeroFloat64)
}

// genPositiveInt generates positive int values
func genPositiveInt() gopter.Gen {
	return gen.IntRange(0, math.MaxInt)
}

// genUTF8String generates valid UTF-8 strings
func genUTF8String() gopter.Gen {
	return gen.AnyString()
}

// genDuration generates time.Duration values
func genDuration() gopter.Gen {
	return gen.Int64().Map(func(n int64) time.Duration {
		return time.Duration(n)
	})
}

// genComplex64 generates complex64 values
func genComplex64() gopter.Gen {
	return gopter.CombineGens(genFloat32(), genFloat32()).Map(func(vals []interface{}) complex64 {
		return complex(vals[0].(float32), vals[1].(float32))
	})
}

// genComplex128 generates complex128 values
func genComplex128() gopter.Gen {
	return gopter.CombineGens(genFloat64(), genFloat64()).Map(func(vals []interface{}) complex128 {
		return complex(vals[0].(float64), vals[1].(float64))
	})
}

// genBigInt generates *big.Int values
func genBigInt() gopter.Gen {
	return gen.Int64().Map(func(n int64) *big.Int {
		return big.NewInt(n)
	})
}

// genBigFloat generates *big.Float values
func genBigFloat() gopter.Gen {
	return genFloat64().Map(func(f float64) *big.Float {
		return big.NewFloat(f)
	})
}

// genBigRat generates *big.Rat values
func genBigRat() gopter.Gen {
	return gopter.CombineGens(gen.Int64(), gen.Int64Range(1, math.MaxInt64)).Map(func(vals []interface{}) *big.Rat {
		return big.NewRat(vals[0].(int64), vals[1].(int64))
	})
}

// genNonZeroInt generates non-zero int values
func genNonZeroInt() gopter.Gen {
	return gen.Int().SuchThat(func(n int) bool {
		return n != 0
	})
}

// genInvalidNumericString generates strings that cannot be parsed as numbers
func genInvalidNumericString() gopter.Gen {
	return gen.OneConstOf("abc", "test", "hello", "world", "NaN", "invalid", "12.34.56", "1e1e1")
}

// =============================================================================
// Helper Types for Testing
// =============================================================================

// testStringer implements fmt.Stringer for testing
type testStringer struct {
	value string
}

func (s testStringer) String() string {
	return s.value
}

// testError implements error for testing
type testError struct {
	message string
}

func (e testError) Error() string {
	return e.message
}

// =============================================================================
// Property Test Configuration
// =============================================================================

func getTestParameters() *gopter.TestParameters {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 100
	return params
}

// =============================================================================
// Property 1: Integer Type Round-Trip Consistency
// **Feature: cast-extension, Property 1: 整数类型相互转换一致性**
// **Validates: Requirements 1.1, 1.2**
// =============================================================================

func TestProperty1_IntegerRoundTrip(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// int8 -> int64 -> int8 round-trip
	properties.Property("int8 round-trip via int64", prop.ForAll(
		func(v int8) bool {
			result := cast.ToInt8(cast.ToInt64(v))
			return result == v
		},
		genInt8InRange(),
	))

	// int16 -> int64 -> int16 round-trip
	properties.Property("int16 round-trip via int64", prop.ForAll(
		func(v int16) bool {
			result := cast.ToInt16(cast.ToInt64(v))
			return result == v
		},
		genInt16InRange(),
	))

	// int32 -> int64 -> int32 round-trip
	properties.Property("int32 round-trip via int64", prop.ForAll(
		func(v int32) bool {
			result := cast.ToInt32(cast.ToInt64(v))
			return result == v
		},
		genInt32InRange(),
	))

	// uint8 -> uint64 -> uint8 round-trip
	properties.Property("uint8 round-trip via uint64", prop.ForAll(
		func(v uint8) bool {
			result := cast.ToUint8(cast.ToUint64(v))
			return result == v
		},
		genUint8(),
	))

	// uint16 -> uint64 -> uint16 round-trip
	properties.Property("uint16 round-trip via uint64", prop.ForAll(
		func(v uint16) bool {
			result := cast.ToUint16(cast.ToUint64(v))
			return result == v
		},
		genUint16(),
	))

	// uint32 -> uint64 -> uint32 round-trip
	properties.Property("uint32 round-trip via uint64", prop.ForAll(
		func(v uint32) bool {
			result := cast.ToUint32(cast.ToUint64(v))
			return result == v
		},
		genUint32(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 2: Negative to Unsigned Returns Error
// **Feature: cast-extension, Property 2: 负数转无符号整数返回错误**
// **Validates: Requirements 1.3, 6.3**
// =============================================================================

func TestProperty2_NegativeToUnsignedError(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// Negative int to uint should error
	properties.Property("negative int to uint returns error", prop.ForAll(
		func(v int) bool {
			if v >= 0 {
				return true // skip non-negative values
			}
			_, err := cast.ToUintE(v)
			return err != nil
		},
		gen.Int(),
	))

	// Negative int to uint8 should error
	properties.Property("negative int to uint8 returns error", prop.ForAll(
		func(v int) bool {
			if v >= 0 {
				return true
			}
			_, err := cast.ToUint8E(v)
			return err != nil
		},
		gen.Int(),
	))

	// Negative int to uint16 should error
	properties.Property("negative int to uint16 returns error", prop.ForAll(
		func(v int) bool {
			if v >= 0 {
				return true
			}
			_, err := cast.ToUint16E(v)
			return err != nil
		},
		gen.Int(),
	))

	// Negative int to uint32 should error
	properties.Property("negative int to uint32 returns error", prop.ForAll(
		func(v int) bool {
			if v >= 0 {
				return true
			}
			_, err := cast.ToUint32E(v)
			return err != nil
		},
		gen.Int(),
	))

	// Negative int to uint64 should error
	properties.Property("negative int to uint64 returns error", prop.ForAll(
		func(v int) bool {
			if v >= 0 {
				return true
			}
			_, err := cast.ToUint64E(v)
			return err != nil
		},
		gen.Int(),
	))

	// Negative float64 to uint should error
	properties.Property("negative float64 to uint returns error", prop.ForAll(
		func(v float64) bool {
			if v >= 0 || math.IsNaN(v) || math.IsInf(v, 0) {
				return true
			}
			_, err := cast.ToUintE(v)
			return err != nil
		},
		gen.Float64(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 3: Float to Int Truncation
// **Feature: cast-extension, Property 3: 浮点数与整数转换保持整数部分**
// **Validates: Requirements 1.4**
// =============================================================================

func TestProperty3_FloatToIntTruncation(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// float64 to int truncates decimal part
	properties.Property("float64 to int truncates", prop.ForAll(
		func(f float64) bool {
			if math.IsNaN(f) || math.IsInf(f, 0) {
				return true
			}
			if f > float64(math.MaxInt64) || f < float64(math.MinInt64) {
				return true // skip overflow cases
			}
			result := cast.ToInt64(f)
			expected := int64(f)
			return result == expected
		},
		gen.Float64(),
	))

	// float32 to int truncates decimal part
	properties.Property("float32 to int truncates", prop.ForAll(
		func(f float32) bool {
			if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
				return true
			}
			if float64(f) > float64(math.MaxInt32) || float64(f) < float64(math.MinInt32) {
				return true
			}
			result := cast.ToInt32(f)
			expected := int32(f)
			return result == expected
		},
		gen.Float32(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 4: Big Number Round-Trip
// **Feature: cast-extension, Property 4: 大数类型与基础数值类型转换一致性**
// **Validates: Requirements 1.5**
// =============================================================================

func TestProperty4_BigNumberRoundTrip(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// int64 -> *big.Int -> int64 round-trip
	properties.Property("int64 to big.Int round-trip", prop.ForAll(
		func(v int64) bool {
			bigInt := cast.ToBigInt(v)
			result := cast.ToInt64(bigInt)
			return result == v
		},
		gen.Int64(),
	))

	// float64 -> *big.Float -> float64 round-trip (within precision)
	properties.Property("float64 to big.Float round-trip", prop.ForAll(
		func(v float64) bool {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return true
			}
			bigFloat := cast.ToBigFloat(v)
			result := cast.ToFloat64(bigFloat)
			// Allow small precision differences
			diff := math.Abs(result - v)
			return diff < 1e-10 || diff/math.Abs(v) < 1e-10
		},
		genFloat64(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 5: Complex to Real Takes Real Part
// **Feature: cast-extension, Property 5: 复数转实数取实部**
// **Validates: Requirements 1.6**
// =============================================================================

func TestProperty5_ComplexToRealTakesRealPart(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// complex64 to float32 takes real part
	properties.Property("complex64 to float32 takes real part", prop.ForAll(
		func(r float32, i float32) bool {
			if math.IsNaN(float64(r)) || math.IsInf(float64(r), 0) {
				return true
			}
			c := complex(r, i)
			result := cast.ToFloat32(c)
			return result == r
		},
		genFloat32(),
		genFloat32(),
	))

	// complex128 to float64 takes real part
	properties.Property("complex128 to float64 takes real part", prop.ForAll(
		func(r float64, i float64) bool {
			if math.IsNaN(r) || math.IsInf(r, 0) {
				return true
			}
			c := complex(r, i)
			result := cast.ToFloat64(c)
			return result == r
		},
		genFloat64(),
		genFloat64(),
	))

	// complex64 to int takes real part (truncated)
	properties.Property("complex64 to int takes real part truncated", prop.ForAll(
		func(r float32, i float32) bool {
			if math.IsNaN(float64(r)) || math.IsInf(float64(r), 0) {
				return true
			}
			if float64(r) > float64(math.MaxInt32) || float64(r) < float64(math.MinInt32) {
				return true
			}
			c := complex(r, i)
			result := cast.ToInt32(c)
			expected := int32(r)
			return result == expected
		},
		genFloat32(),
		genFloat32(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 6: String and Bytes Round-Trip
// **Feature: cast-extension, Property 6: 字符串与字节切片往返一致性**
// **Validates: Requirements 2.1**
// =============================================================================

func TestProperty6_StringBytesRoundTrip(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// string -> []byte -> string round-trip
	properties.Property("string to bytes round-trip", prop.ForAll(
		func(s string) bool {
			bytes := cast.ToBytes(s)
			result := cast.ToString(bytes)
			return result == s
		},
		gen.AnyString(),
	))

	// []byte -> string -> []byte round-trip
	properties.Property("bytes to string round-trip", prop.ForAll(
		func(s string) bool {
			bytes := []byte(s)
			str := cast.ToString(bytes)
			result := cast.ToBytes(str)
			if len(bytes) != len(result) {
				return false
			}
			for i := range bytes {
				if bytes[i] != result[i] {
					return false
				}
			}
			return true
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 7: Any Type to String Succeeds
// **Feature: cast-extension, Property 7: 任意类型转字符串类成功**
// **Validates: Requirements 2.2, 2.3, 2.4, 2.5**
// =============================================================================

func TestProperty7_AnyTypeToStringSucceeds(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// int to string succeeds
	properties.Property("int to string succeeds", prop.ForAll(
		func(v int) bool {
			_, err := cast.ToStringE(v)
			return err == nil
		},
		gen.Int(),
	))

	// float64 to string succeeds
	properties.Property("float64 to string succeeds", prop.ForAll(
		func(v float64) bool {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return true
			}
			_, err := cast.ToStringE(v)
			return err == nil
		},
		genFloat64(),
	))

	// bool to string succeeds
	properties.Property("bool to string succeeds", prop.ForAll(
		func(v bool) bool {
			_, err := cast.ToStringE(v)
			return err == nil
		},
		gen.Bool(),
	))

	// int to bytes succeeds
	properties.Property("int to bytes succeeds", prop.ForAll(
		func(v int) bool {
			_, err := cast.ToBytesE(v)
			return err == nil
		},
		gen.Int(),
	))

	// int to stringer succeeds
	properties.Property("int to stringer succeeds", prop.ForAll(
		func(v int) bool {
			_, err := cast.ToStringerE(v)
			return err == nil
		},
		gen.Int(),
	))

	// int to error succeeds
	properties.Property("int to error succeeds", prop.ForAll(
		func(v int) bool {
			_, err := cast.ToErrorE(v)
			return err == nil
		},
		gen.Int(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 8: Stringer/Error Input Uses String Method
// **Feature: cast-extension, Property 8: fmt.Stringer 和 error 输入使用字符串方法**
// **Validates: Requirements 2.6, 2.7**
// =============================================================================

func TestProperty8_StringerErrorInputUsesStringMethod(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// fmt.Stringer input uses String() for numeric conversion
	properties.Property("stringer input uses String() for int conversion", prop.ForAll(
		func(v int64) bool {
			s := testStringer{value: fmt.Sprintf("%d", v)}
			result, err := cast.ToInt64E(s)
			if err != nil {
				return false
			}
			return result == v
		},
		gen.Int64(),
	))

	// error input uses Error() for numeric conversion
	properties.Property("error input uses Error() for int conversion", prop.ForAll(
		func(v int64) bool {
			e := errors.New(fmt.Sprintf("%d", v))
			result, err := cast.ToInt64E(e)
			if err != nil {
				return false
			}
			return result == v
		},
		gen.Int64(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 9: Duration and Numeric Round-Trip
// **Feature: cast-extension, Property 9: Duration 与数值类型往返一致性**
// **Validates: Requirements 3.1**
// =============================================================================

func TestProperty9_DurationNumericRoundTrip(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// time.Duration -> int64 -> time.Duration round-trip
	properties.Property("duration to int64 round-trip", prop.ForAll(
		func(d int64) bool {
			duration := time.Duration(d)
			n := cast.ToInt64(duration)
			result := cast.ToDuration(n)
			return result == duration
		},
		gen.Int64(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 10: Numeric to Bool Rule
// **Feature: cast-extension, Property 10: 数值转布尔规则**
// **Validates: Requirements 4.1, 4.4**
// =============================================================================

func TestProperty10_NumericToBoolRule(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// int to bool: non-zero is true, zero is false
	properties.Property("int to bool: non-zero is true", prop.ForAll(
		func(v int) bool {
			result := cast.ToBool(v)
			if v == 0 {
				return result == false
			}
			return result == true
		},
		gen.Int(),
	))

	// float64 to bool: non-zero is true, zero is false
	properties.Property("float64 to bool: non-zero is true", prop.ForAll(
		func(v float64) bool {
			if math.IsNaN(v) {
				return true
			}
			result := cast.ToBool(v)
			if v == 0 {
				return result == false
			}
			return result == true
		},
		gen.Float64(),
	))

	// complex to bool: real or imag non-zero is true
	properties.Property("complex64 to bool: non-zero is true", prop.ForAll(
		func(r float32, i float32) bool {
			if math.IsNaN(float64(r)) || math.IsNaN(float64(i)) {
				return true
			}
			c := complex(r, i)
			result := cast.ToBool(c)
			if r == 0 && i == 0 {
				return result == false
			}
			return result == true
		},
		gen.Float32(),
		gen.Float32(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 11: Bool to Numeric Rule
// **Feature: cast-extension, Property 11: 布尔转数值规则**
// **Validates: Requirements 4.2**
// =============================================================================

func TestProperty11_BoolToNumericRule(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// bool to int: true is 1, false is 0
	properties.Property("bool to int: true is 1, false is 0", prop.ForAll(
		func(b bool) bool {
			result := cast.ToInt(b)
			if b {
				return result == 1
			}
			return result == 0
		},
		gen.Bool(),
	))

	// bool to float64: true is 1, false is 0
	properties.Property("bool to float64: true is 1, false is 0", prop.ForAll(
		func(b bool) bool {
			result := cast.ToFloat64(b)
			if b {
				return result == 1.0
			}
			return result == 0.0
		},
		gen.Bool(),
	))

	// bool to uint: true is 1, false is 0
	properties.Property("bool to uint: true is 1, false is 0", prop.ForAll(
		func(b bool) bool {
			result := cast.ToUint(b)
			if b {
				return result == 1
			}
			return result == 0
		},
		gen.Bool(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 12: Pointer Dereference
// **Feature: cast-extension, Property 12: 指针解引用正确性**
// **Validates: Requirements 5.2**
// =============================================================================

func TestProperty12_PointerDereference(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// *int to int dereferences correctly
	properties.Property("*int to int dereferences", prop.ForAll(
		func(v int) bool {
			ptr := &v
			result := cast.ToInt(ptr)
			return result == v
		},
		gen.Int(),
	))

	// *string to string dereferences correctly
	properties.Property("*string to string dereferences", prop.ForAll(
		func(s string) bool {
			ptr := &s
			result := cast.ToString(ptr)
			return result == s
		},
		gen.AnyString(),
	))

	// **int to int dereferences correctly (double pointer)
	properties.Property("**int to int dereferences", prop.ForAll(
		func(v int) bool {
			ptr := &v
			pptr := &ptr
			result := cast.ToInt(pptr)
			return result == v
		},
		gen.Int(),
	))

	properties.TestingRun(t)
}

// =============================================================================
// Property 13: Nil Input Returns Zero Value
// **Feature: cast-extension, Property 13: nil 输入返回零值**
// **Validates: Requirements 6.1**
// =============================================================================

func TestProperty13_NilInputReturnsZeroValue(t *testing.T) {
	// Test nil input for all conversion functions
	t.Run("nil to int", func(t *testing.T) {
		result, err := cast.ToIntE(nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("nil to int64", func(t *testing.T) {
		result, err := cast.ToInt64E(nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("nil to uint", func(t *testing.T) {
		result, err := cast.ToUintE(nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("nil to float64", func(t *testing.T) {
		result, err := cast.ToFloat64E(nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != 0 {
			t.Errorf("expected 0, got %f", result)
		}
	})

	t.Run("nil to string", func(t *testing.T) {
		result, err := cast.ToStringE(nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	t.Run("nil to bool", func(t *testing.T) {
		result, err := cast.ToBoolE(nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != false {
			t.Errorf("expected false, got %v", result)
		}
	})

	t.Run("nil to duration", func(t *testing.T) {
		result, err := cast.ToDurationE(nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != 0 {
			t.Errorf("expected 0, got %v", result)
		}
	})
}

// =============================================================================
// Property 14: Nil Big Number Pointer Returns Error
// **Feature: cast-extension, Property 14: nil 大数指针返回错误**
// **Validates: Requirements 6.4**
// =============================================================================

func TestProperty14_NilBigNumberPointerReturnsError(t *testing.T) {
	t.Run("nil *big.Int to int", func(t *testing.T) {
		var nilBigInt *big.Int
		_, err := cast.ToIntE(nilBigInt)
		if err == nil {
			t.Error("expected error for nil *big.Int")
		}
	})

	t.Run("nil *big.Float to float64", func(t *testing.T) {
		var nilBigFloat *big.Float
		_, err := cast.ToFloat64E(nilBigFloat)
		if err == nil {
			t.Error("expected error for nil *big.Float")
		}
	})

	t.Run("nil *big.Rat to float64", func(t *testing.T) {
		var nilBigRat *big.Rat
		_, err := cast.ToFloat64E(nilBigRat)
		if err == nil {
			t.Error("expected error for nil *big.Rat")
		}
	})
}

// =============================================================================
// Property 15: Invalid String Returns Error
// **Feature: cast-extension, Property 15: 无效字符串返回错误**
// **Validates: Requirements 6.5**
// =============================================================================

func TestProperty15_InvalidStringReturnsError(t *testing.T) {
	params := getTestParameters()
	properties := gopter.NewProperties(params)

	// Invalid string to int returns error
	properties.Property("invalid string to int returns error", prop.ForAll(
		func(s string) bool {
			_, err := cast.ToIntE(s)
			return err != nil
		},
		genInvalidNumericString(),
	))

	// Invalid string to float64 returns error
	properties.Property("invalid string to float64 returns error", prop.ForAll(
		func(s string) bool {
			_, err := cast.ToFloat64E(s)
			return err != nil
		},
		genInvalidNumericString(),
	))

	// Invalid string to uint returns error
	properties.Property("invalid string to uint returns error", prop.ForAll(
		func(s string) bool {
			_, err := cast.ToUintE(s)
			return err != nil
		},
		genInvalidNumericString(),
	))

	properties.TestingRun(t)
}
