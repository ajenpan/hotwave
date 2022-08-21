
$release_dir="bin"

If (!(test-path $release_dir)){
    md $release_dir
}

# go env -w GOOS="windows"
# go build -o $release_dir/gateway.exe ./apps/gateway
# go build -o $release_dir/battle.exe ./apps/battle
# go build -o $release_dir/auth.exe ./apps/auth


# build for linux
go env -w GOOS="linux"
go build -o $release_dir/gateway ./apps/gateway


go env -w GOOS="windows"
