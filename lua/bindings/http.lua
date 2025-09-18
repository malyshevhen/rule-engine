---@meta

---@module 'http' HTTP module that provides functions to make HTTP requests
local http

---@class http.Response HTTP response
---@field status integer HTTP status code
---@field body string? HTTP response body

--- Get makes an HTTP GET request
---
---@param url string HTTP URL to make the request to
---@param headers table<string, string> HTTP headers to send with the request
---@return http.Response? HTTP response
---@return string? Error message
function http.get(url, headers) end

--- Post makes an HTTP POST request
---
---@param url string HTTP URL to make the request to
---@param headers table<string, string> HTTP headers to send with the request
---@param body string HTTP body to send with the request
---@return http.Response? HTTP response
---@return string? Error message
function http.post(url, headers, body) end

--- Delete makes an HTTP DELETE request
---
---@param url string HTTP URL to make the request to
---@param headers table<string, string> HTTP headers to send with the request
---@return http.Response? HTTP response
---@return string? Error message
function http.delete(url, headers) end

--- Put makes an HTTP PUT request
---
---@param url string HTTP URL to make the request to
---@param headers table<string, string> HTTP headers to send with the request
---@param body string HTTP body to send with the request
---@return http.Response? HTTP response
---@return string? Error message
function http.put(url, headers, body) end

--- Patch makes an HTTP PATCH request
---
---@param url string HTTP URL to make the request to
---@param headers table<string, string> HTTP headers to send with the request
---@param body string HTTP body to send with the request
---@return http.Response? HTTP response
---@return string? Error message
function http.patch(url, headers, body) end

return http
