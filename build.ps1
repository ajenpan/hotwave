
$release_dir="bin"

If (!(test-path $release_dir)){
    md $release_dir
}

$GitCommit=git log -1 --format="%ct"

go env -w GOOS="windows"
go build -o $release_dir/bfmanager.exe -mod=vendor -ldflags "-X main.GitCommit=$GitCommit" ./cmd

go env -w GOOS="linux"
go build -o account .\servers\account\cmd\ 