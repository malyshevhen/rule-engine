local http = require 'http'
local log = require 'logger'

log.info 'Sending HTTP POST request'

local response, err = http.post('http://hoverfly:8888/api/webhook', {}, '')
if err then
  log.error(err)
  return false
end

if not response then
  log.error 'No response received'
  return false
end

return response.status == 200 and response.body == '{"status": "success"}'
