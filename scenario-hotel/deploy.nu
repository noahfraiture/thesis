cd tools/red; gcc robot.c -o robot; cd -
cp tools/red/robot ./nix/hosts/dashboard/robot
cd tools/hikvision/server ; go build ; cd -
cp tools/hikvision/server/main ./nix/hosts/camera/hikvision
let dir = gum choose ...(ls vms | get name)
let vm = $dir | parse "vms/{name}-{port}" | first
let port = $vm | get port
let name = $vm | get name
quickemu --ssh-port $port --vm $"($dir)/nixos-24.11-gnome.conf"
try {
  gum confirm "Do you want to build ?"
  ssh -t env-admin@localhost -p $port 'echo env-admin | sudo -S rm -rf /home/env-admin/nixos'
  scp -r -P $port /home/noah/Projects/thesis/scenario-hotel/nix env-admin@localhost:/home/env-admin/nixos
  ssh -t env-admin@localhost -p $port 'echo env-admin | sudo -S cp /etc/nixos/hardware-configuration.nix /home/env-admin/nixos/hosts/'
  ssh -t env-admin@localhost -p $port 'echo env-admin | sudo -S chmod 0700 /home/env-admin/nixos'

  # ssh -t env-admin@localhost -p $port $'echo env-admin | sudo -S nix-collect-garbage -d'
  ssh -t env-admin@localhost -p $port $'echo env-admin | sudo -S nixos-rebuild switch --flake /home/env-admin/nixos/#($name)'
}
print $"Start vm with port (ansi red)($port)(ansi reset) and config (ansi red)($name)(ansi reset)"
