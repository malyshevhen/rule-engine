local log = require 'logger'
local time = require 'time'

log.info 'Starting script...'

local function getCurrentTime()
  return time.now()
end

log.info('Current time is: ' .. getCurrentTime())

local currentTime = getCurrentTime()
log.info('Current time is: ' .. currentTime)

if currentTime == nil or currentTime == '' then
  log.error 'Current time is not set'
  return false
end

return true
