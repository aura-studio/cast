# 需求文档

## 简介

本文档定义了 Go 语言类型转换库 `github.com/aura-studio/cast` 的类型转换完整性检查需求。目标是确保库支持的所有类型之间能够进行 N×N 的完整转换。

## 术语表

- **Cast_Library**: 类型转换库，提供各种类型之间的安全转换功能
- **Source_Type**: 源类型，作为转换输入的类型
- **Target_Type**: 目标类型，转换输出的类型
- **N×N_Coverage**: 完整覆盖，指所有源类型都能转换到所有目标类型
- **E_Suffix_Function**: 带 E 后缀的函数，返回 (结果, error) 元组，用于错误处理

## 当前支持的类型

### 源类型（输入类型）
- 有符号整数：int, int8, int16, int32, int64
- 无符号整数：uint, uint8, uint16, uint32, uint64
- 浮点数：float32, float64
- 大数类型：*big.Int, *big.Float, *big.Rat
- 复数类型：complex64, complex128
- 布尔类型：bool
- 时间类型：time.Duration, time.Location, *time.Location
- 字符串类型：string, []byte, fmt.Stringer, error
- 空值：nil

### 目标类型（输出类型）
- 有符号整数：int, int8, int16, int32, int64
- 无符号整数：uint, uint8, uint16, uint32, uint64
- 浮点数：float32, float64
- 大数类型：*big.Int, *big.Float, *big.Rat
- 复数类型：complex64, complex128
- 布尔类型：bool
- 时间类型：time.Duration, *time.Location
- 字符串类型：string, []byte, fmt.Stringer, error

## 需求

### 需求 1：数值类型之间的完整转换

**用户故事：** 作为开发者，我希望所有数值类型之间都能相互转换，以便在处理不同精度和范围的数值时能够灵活使用。

#### 验收标准

1. FOR ALL 有符号整数类型（int, int8, int16, int32, int64），Cast_Library SHALL 支持相互转换
2. FOR ALL 无符号整数类型（uint, uint8, uint16, uint32, uint64），Cast_Library SHALL 支持相互转换
3. FOR ALL 有符号整数和无符号整数类型，Cast_Library SHALL 支持相互转换（负数转无符号时返回错误）
4. FOR ALL 浮点数类型（float32, float64），Cast_Library SHALL 支持与整数类型的相互转换
5. FOR ALL 大数类型（*big.Int, *big.Float, *big.Rat），Cast_Library SHALL 支持与基础数值类型的相互转换
6. FOR ALL 复数类型（complex64, complex128），Cast_Library SHALL 支持与实数类型的相互转换（取实部）

### 需求 2：字符串类型之间的完整转换

**用户故事：** 作为开发者，我希望字符串相关类型之间都能相互转换，以便在处理文本数据时能够灵活选择表示方式。

#### 验收标准

1. THE Cast_Library SHALL 支持 string 与 []byte 的相互转换
2. THE Cast_Library SHALL 支持任意类型转换为 string
3. THE Cast_Library SHALL 支持任意类型转换为 []byte
4. THE Cast_Library SHALL 支持任意类型转换为 fmt.Stringer
5. THE Cast_Library SHALL 支持任意类型转换为 error
6. WHEN 输入为 fmt.Stringer 类型 THEN Cast_Library SHALL 使用 String() 方法获取字符串后进行转换
7. WHEN 输入为 error 类型 THEN Cast_Library SHALL 使用 Error() 方法获取字符串后进行转换

### 需求 3：时间类型之间的完整转换

**用户故事：** 作为开发者，我希望时间相关类型之间都能相互转换，以便在处理时间和时区数据时能够灵活使用。

#### 验收标准

1. THE Cast_Library SHALL 支持 time.Duration 与数值类型的相互转换
2. THE Cast_Library SHALL 支持 time.Location 与 *time.Location 的相互转换
3. THE Cast_Library SHALL 支持 time.Duration 与 *time.Location 的相互转换（通过时区偏移量）
4. THE Cast_Library SHALL 支持字符串与 time.Duration 的相互转换
5. THE Cast_Library SHALL 支持字符串与 *time.Location 的相互转换

### 需求 4：布尔类型的完整转换

**用户故事：** 作为开发者，我希望布尔类型能够与其他类型相互转换，以便在处理条件判断和标志位时能够灵活使用。

#### 验收标准

1. FOR ALL 数值类型，Cast_Library SHALL 支持转换为 bool（非零为 true，零为 false）
2. FOR ALL 数值类型，Cast_Library SHALL 支持从 bool 转换（true 为 1，false 为 0）
3. THE Cast_Library SHALL 支持字符串与 bool 的相互转换（"true"/"false" 等）
4. FOR ALL 复数类型，Cast_Library SHALL 支持转换为 bool（实部或虚部非零为 true）

### 需求 5：API 一致性

**用户故事：** 作为开发者，我希望所有转换函数保持一致的 API 风格，以便降低学习成本并保持代码风格统一。

#### 验收标准

1. THE Cast_Library SHALL 为每个目标类型提供两个版本的转换函数：简单版本（直接返回结果）和带 E 后缀版本（返回结果和错误）
2. THE Cast_Library SHALL 在所有转换函数中使用 `indirectToStringerOrError` 处理指针解引用
3. THE Cast_Library SHALL 在所有转换函数中使用 type switch 处理各种输入类型
4. THE Cast_Library SHALL 保持一致的错误消息格式："unable to cast %#v of type %T to [目标类型]"

### 需求 6：边界情况处理

**用户故事：** 作为开发者，我希望转换函数能够正确处理各种边界情况，以便在生产环境中安全使用。

#### 验收标准

1. WHEN 输入为 nil THEN 所有转换函数 SHALL 返回对应类型的零值
2. WHEN 数值溢出目标类型范围 THEN 转换函数 SHALL 进行截断（与 Go 语言标准行为一致）
3. WHEN 负数转换为无符号类型 THEN 转换函数 SHALL 返回错误
4. WHEN *big.Int、*big.Float 或 *big.Rat 为 nil 指针 THEN 转换函数 SHALL 返回错误
5. WHEN 字符串无法解析为目标数值类型 THEN 转换函数 SHALL 返回错误
