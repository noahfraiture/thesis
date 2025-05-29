let dir = gum choose ...(ls vms | get name)
let vm = $dir | parse "vms/{name}-{port}" | first
let port = $vm | get port
let name = $vm | get name
print $"Connect vm with port (ansi red)($port)(ansi reset) and config (ansi red)($name)(ansi reset)"
ssh env-admin@localhost -p $port
