set tabstop=8
set softtabstop=0 noexpandtab
set shiftwidth=8

autocmd FileType go setlocal shiftwidth=8 softtabstop=0 noexpandtab tabstop=8
autocmd FileType yaml setlocal ts=2 sts=2 sw=2 expandtab 
let g:go_fmt_command = "goimports"
