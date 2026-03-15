" blockfont project-specific Vim/Neovim configuration
" Load with :source .exrc or use exrc option

" Go-specific settings
setlocal tabstop=4
setlocal shiftwidth=4
setlocal noexpandtab

" Which-key mappings under <leader>p for project commands
lua << EOF
local ok, wk = pcall(require, "which-key")
if ok then
  wk.add({
    -- Main project group
    { "<leader>p", group = "Project (blockfont)" },

    -- Build & Run
    { "<leader>pb", "<cmd>!make build<CR>", desc = "Build library" },
    { "<leader>pd", "<cmd>!make demo<CR>", desc = "Run demo" },
    { "<leader>pt", "<cmd>!make test<CR>", desc = "Run tests" },
    { "<leader>pT", "<cmd>!make test-coverage<CR>", desc = "Test coverage" },

    -- Examples subgroup
    { "<leader>pe", group = "Examples" },
    { "<leader>pes", "<cmd>!make run-simple<CR>", desc = "Simple example" },
    { "<leader>pea", "<cmd>!make run-animated<CR>", desc = "Animated example" },
    { "<leader>pee", "<cmd>!make run-editor<CR>", desc = "Editor example" },

    -- Code Quality
    { "<leader>pf", "<cmd>!make fmt<CR>", desc = "Format code" },
    { "<leader>pl", "<cmd>!make lint<CR>", desc = "Lint code" },

    -- Documentation subgroup
    { "<leader>pD", group = "Documentation" },
    { "<leader>pDd", "<cmd>!make docs-dev &<CR>", desc = "Start docs server" },
    { "<leader>pDb", "<cmd>!make docs-build<CR>", desc = "Build docs" },
    { "<leader>pDc", "<cmd>!make docs-clean<CR>", desc = "Clean docs" },
    { "<leader>pDo", "<cmd>!make docs-open<CR>", desc = "Open docs" },

    -- Demo Recording subgroup
    { "<leader>pr", group = "Recording" },
    { "<leader>prr", "<cmd>terminal nix run .#demo-record<CR>", desc = "Record demo" },
    { "<leader>prp", "<cmd>terminal nix run .#demo-play<CR>", desc = "Play demo" },

    -- Release
    { "<leader>pR", "<cmd>terminal nix run .#release<CR>", desc = "Release" },

    -- Nix subgroup
    { "<leader>pn", group = "Nix" },
    { "<leader>pnb", "<cmd>!nix build<CR>", desc = "Nix build" },
    { "<leader>pnr", "<cmd>!nix run<CR>", desc = "Nix run" },
    { "<leader>pnd", "<cmd>!nix run .#docs-serve<CR>", desc = "Nix docs serve" },
    { "<leader>pnt", "<cmd>!nix run .#test<CR>", desc = "Nix test" },

    -- Version & Dependencies
    { "<leader>pv", "<cmd>!make version<CR>", desc = "Show version" },
    { "<leader>pu", "<cmd>!make deps<CR>", desc = "Update deps" },
    { "<leader>pV", "<cmd>!make vendor<CR>", desc = "Vendor deps" },
  })
end
EOF

" Fallback keybindings for regular Vim (without which-key)
if !has('nvim')
  nnoremap <leader>pb :!make build<CR>
  nnoremap <leader>pd :!make demo<CR>
  nnoremap <leader>pt :!make test<CR>
  nnoremap <leader>pf :!make fmt<CR>
  nnoremap <leader>pl :!make lint<CR>
  nnoremap <leader>pDd :!make docs-dev &<CR>
  nnoremap <leader>pv :!make version<CR>
endif
