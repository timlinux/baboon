-- Baboon project-specific Neovim configuration
-- This file is automatically loaded when opening Neovim in the project directory

-- Set Go-specific options
vim.opt_local.tabstop = 4
vim.opt_local.shiftwidth = 4
vim.opt_local.expandtab = false

-- Project root detection
vim.g.baboon_root = vim.fn.getcwd()

-- Custom commands for this project
vim.api.nvim_create_user_command('BaboonBuild', '!make build', {})
vim.api.nvim_create_user_command('BaboonRun', '!make run', {})
vim.api.nvim_create_user_command('BaboonTest', '!make test', {})
vim.api.nvim_create_user_command('BaboonFmt', '!make fmt', {})
vim.api.nvim_create_user_command('BaboonLint', '!make lint', {})
vim.api.nvim_create_user_command('BaboonDocsDev', '!make docs-dev &', {})
vim.api.nvim_create_user_command('BaboonDocsBuild', '!make docs-build', {})
vim.api.nvim_create_user_command('BaboonDocsOpen', '!make docs-open', {})
vim.api.nvim_create_user_command('BaboonWebDev', '!make web-dev', {})

-- Which-key integration for project keybindings
local ok, wk = pcall(require, "which-key")
if ok then
  wk.register({
    p = {
      name = "Project (Baboon)",
      -- Build & Run
      b = { "<cmd>!make build<cr>", "Build binary" },
      r = { "<cmd>!make run<cr>", "Run application" },
      R = { "<cmd>!make run-p<cr>", "Run with punctuation" },
      t = { "<cmd>!make test<cr>", "Run tests" },

      -- Server/Client
      s = { "<cmd>!make server<cr>", "Start server only" },
      c = { "<cmd>!make client<cr>", "Start client only" },
      S = { "<cmd>!make start-backend<cr>", "Start backend (background)" },
      K = { "<cmd>!make stop-backend<cr>", "Stop backend" },

      -- Code Quality
      f = { "<cmd>!make fmt<cr>", "Format code" },
      l = { "<cmd>!make lint<cr>", "Lint code" },

      -- Documentation (Hugo)
      d = {
        name = "Documentation",
        d = { "<cmd>!make docs-dev &<cr>", "Start Hugo dev server" },
        b = { "<cmd>!make docs-build<cr>", "Build documentation" },
        c = { "<cmd>!make docs-clean<cr>", "Clean documentation" },
        o = { "<cmd>!make docs-open<cr>", "Open docs in browser" },
        n = { "<cmd>terminal make docs-new<cr>", "New documentation page" },
      },

      -- Web Frontend
      w = {
        name = "Web Frontend",
        i = { "<cmd>!make web-install<cr>", "Install web dependencies" },
        d = { "<cmd>!make web-dev<cr>", "Start web dev server" },
        b = { "<cmd>!make web-build<cr>", "Build web for production" },
        s = { "<cmd>!make web-start<cr>", "Start backend + web" },
      },

      -- Nix
      n = {
        name = "Nix",
        b = { "<cmd>!nix build<cr>", "Nix build" },
        r = { "<cmd>!nix run<cr>", "Nix run" },
        d = { "<cmd>!nix run .#docs-serve<cr>", "Nix docs serve" },
      },
    },
  }, { prefix = "<leader>" })
end

-- Autocommands for this project
local baboon_group = vim.api.nvim_create_augroup("BaboonProject", { clear = true })

-- Format Go files on save
vim.api.nvim_create_autocmd("BufWritePre", {
  group = baboon_group,
  pattern = "*.go",
  callback = function()
    vim.lsp.buf.format({ async = false })
  end,
})

-- Set filetype for Hugo templates
vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
  group = baboon_group,
  pattern = vim.g.baboon_root .. "/hugo/themes/baboon/layouts/**/*.html",
  callback = function()
    vim.bo.filetype = "html.gotmpl"
  end,
})

-- Print confirmation
vim.notify("Baboon project configuration loaded", vim.log.levels.INFO)
