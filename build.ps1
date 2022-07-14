
$release_dir="bin"

If (!(test-path $release_dir)){
    md $release_dir
}

$GitCommit=git log -1 --format="%ct"

#go env -w GOOS="windows"
#go env -w GOOS="linux"

