log_message('info', 'Sending HTTP POST request')

local response, err = http_post('http://hoverfly:8888/api/webhook', {}, '')
if err then
  return false
end

return response.status == 200 and response.body == '{"status": "success"}'
