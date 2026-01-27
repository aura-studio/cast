# Cast

Go 语言类型转换工具库，支持各种基础类型之间的安全转换。

## 安装

```bash
go get github.com/aura-studio/cast
```

## 功能

### 数值类型转换

```go
cast.ToInt(v)       // int
cast.ToInt8(v)      // int8
cast.ToInt16(v)     // int16
cast.ToInt32(v)     // int32
cast.ToInt64(v)     // int64
cast.ToUint(v)      // uint
cast.ToUint8(v)     // uint8
cast.ToUint16(v)    // uint16
cast.ToUint32(v)    // uint32
cast.ToUint64(v)    // uint64
cast.ToFloat32(v)   // float32
cast.ToFloat64(v)   // float64
cast.ToBool(v)      // bool
```

### 大数类型转换

```go
cast.ToBigInt(v)    // *big.Int
cast.ToBigFloat(v)  // *big.Float
cast.ToBigRat(v)    // *big.Rat
cast.ToComplex64(v) // complex64
cast.ToComplex128(v)// complex128
```

### 字符串类型转换

```go
cast.ToString(v)    // string
cast.ToBytes(v)     // []byte
cast.ToStringer(v)  // fmt.Stringer
cast.ToError(v)     // error
```

### 时间类型转换

```go
cast.ToDuration(v)  // time.Duration
cast.ToTimeZone(v)  // *time.Location
```

## 错误处理

每个转换函数都有对应的 `E` 后缀版本，返回错误信息：

```go
val, err := cast.ToIntE("123")
if err != nil {
    // 处理错误
}
```

## 支持的输入类型

- 基础类型：int, int8-64, uint, uint8-64, float32, float64, bool
- 大数类型：*big.Int, *big.Float, *big.Rat, complex64, complex128
- 字符串类型：string, []byte, fmt.Stringer, error
- 时间类型：time.Duration, time.Location

## License

MIT
