local http = require 'http'
local log = require 'logger'
local time = require 'time'

log.info('Sending HTTP POST request. Current time: ' .. time.now(time.RFC3339))

local response, err = http.post('http://hoverfly:8888/api/webhook', {}, '')
if err then
  log.error(err)
  return false
end

if not response then
  log.error('No response received. Current time: ' .. time.now(time.RFC3339))
  return false
end

log.info(
  'Response status: ' .. tostring(response.status) .. '. Current time: ' .. time.now(time.RFC3339)
)

return response.status == 200 and response.body == '{"status": "success"}'
