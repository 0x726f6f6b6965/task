package helper

var (
	AddTask string = `
		redis.call("SET", KEYS[1], ARGV[1])
		local op = redis.pcall("ZADD", KEYS[2], ARGV[2], ARGV[2])
		if (op ~= 1) then
			redis.call("DEL", KEYS[1])
			error(op)
		end
		return
	`

	DeleteTask string = `
		local val = redis.call("GET", KEYS[1])
		redis.call("DEL", KEYS[1])
		local op = redis.pcall("ZREM", KEYS[2], ARGV[1])
		if (op ~= 1) then 
			redis.call("SET", KEYS[1], val)
			error(op)
		end
		return
	`
)
