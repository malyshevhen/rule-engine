---@meta

---@module 'time' Time module that provides functions to get the current time
local time

--- Time formats.
---@enum time.Format
local Format = {
  ANSIC = 'Mon Jan _2 15:04:05 2006',
  UnixDate = 'Mon Jan _2 15:04:05 MST 2006',
  RubyDate = 'Mon Jan 02 15:04:05 -0700 2006',
  RFC822 = '02-Jan-06 15:04 MST',
  RFC822Z = '02-Jan-06 15:04 -0700',
  RFC850 = 'Monday, 02-Jan-06 15:04:05 MST',
  RFC1123 = 'Mon, 02 Jan 2006 15:04:05 MST',
  RFC1123Z = 'Mon, 02 Jan 2006 15:04:05 -0700',
  RFC3339 = '2006-01-02T15:04:05Z07:00',
  RFC3339Nano = '2006-01-02T15:04:05.999999999Z07:00',
  Kitchen = '3:04PM',
  Stamp = 'Jan _2 15:04:05',
  StampMilli = 'Jan _2 15:04:05.000',
  StampMicro = 'Jan _2 15:04:05.000000',
  StampNano = 'Jan _2 15:04:05.000000000',
  DateTime = '2006-01-02 15:04:05',
  DateOnly = '2006-01-02',
  TimeOnly = '15:04:05',
}

--- Exported Time formats
---@type time.Format
time.ANSIC = Format.ANSIC
---@type time.Format
time.UnixDate = Format.UnixDate
---@type time.Format
time.RubyDate = Format.RubyDate
---@type time.Format
time.RFC822 = Format.RFC822
---@type time.Format
time.RFC822Z = Format.RFC822Z
---@type time.Format
time.RFC850 = Format.RFC850
---@type time.Format
time.RFC1123 = Format.RFC1123
---@type time.Format
time.RFC1123Z = Format.RFC1123Z
---@type time.Format
time.RFC3339 = Format.RFC3339
---@type time.Format
time.RFC3339Nano = Format.RFC3339Nano
---@type time.Format
time.Kitchen = Format.Kitchen
---@type time.Format
time.Stamp = Format.Stamp
---@type time.Format
time.StampMilli = Format.StampMilli
---@type time.Format
time.StampMicro = Format.StampMicro
---@type time.Format
time.StampNano = Format.StampNano
---@type time.Format
time.DateTime = Format.DateTime
---@type time.Format
time.DateOnly = Format.DateOnly
---@type time.Format
time.TimeOnly = Format.TimeOnly

--- GetCurrentTime returns the current time as a string
---
--- Example:
--- ```
--- local time = require 'time'
---
--- local now = time.now(time.ANSIC)
--- assert(now == 'Mon Jan _2 15:04:05 2006')
---
--- local now = time.now(time.DateTime)
--- assert(now == '2023-01-01 15:04:05')
---
--- now = time.now(time.DateOnly)
--- assert(now == '2023-01-01')
---
--- now = time.now(time.TimeOnly)
--- assert(now == '15:04:05')
---
--- now = time.now()
--- assert(now == '2023-01-01 15:04:05.000000000')
--- ```
--- The default format is the Go time format, which is "2006-01-02 15:04:05.999999999"
---
---@param format time.Format? Formatter to use for the current time
---@return string? time Current time as a string
function time.now(format) end

return time
