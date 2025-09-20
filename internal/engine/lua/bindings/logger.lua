---@meta

---@module 'logger' Logger module that provides functions to log messages
local logger

--- Log a message with the `INFO` level
---
---@param message string Message to log
function logger.info(message) end

--- Log a message with the `DEBUG` level
---
---@param message string Message to log
function logger.debug(message) end

--- Log a message with the `WARN` level
---
---@param message string Message to log
function logger.warn(message) end

--- Log a message with the `ERROR` level
---
---@param message string Message to log
function logger.error(message) end

return logger
