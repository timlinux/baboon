-- blockfont project-specific Neovim configuration
-- This file is automatically loaded when opening Neovim in the project directory

-- Set Go-specific options
vim.opt_local.tabstop = 4
vim.opt_local.shiftwidth = 4
vim.opt_local.expandtab = false

-- Project root detection
vim.g.blockfont_root = vim.fn.getcwd()

-- Custom commands for this project
vim.api.nvim_create_user_command('BlockfontBuild', '!make build', {})
vim.api.nvim_create_user_command('BlockfontTest', '!make test', {})
vim.api.nvim_create_user_command('BlockfontDemo', '!make demo', {})
vim.api.nvim_create_user_command('BlockfontFmt', '!make fmt', {})
vim.api.nvim_create_user_command('BlockfontLint', '!make lint', {})
vim.api.nvim_create_user_command('BlockfontDocsDev', '!make docs-dev &', {})
vim.api.nvim_create_user_command('BlockfontDocsBuild', '!make docs-build', {})
vim.api.nvim_create_user_command('BlockfontDocsOpen', '!make docs-open', {})
vim.api.nvim_create_user_command('BlockfontDemoRecord', '!nix run .#demo-record', {})
vim.api.nvim_create_user_command('BlockfontDemoPlay', '!nix run .#demo-play', {})

-- Which-key integration for project keybindings
local ok, wk = pcall(require, "which-key")
if ok then
  wk.register({
    p = {
      name = "Project (blockfont)",
      -- Build & Run
      b = { "<cmd>!make build<cr>", "Build library" },
      d = { "<cmd>!make demo<cr>", "Run demo application" },
      t = { "<cmd>!make test<cr>", "Run tests" },
      T = { "<cmd>!make test-coverage<cr>", "Run tests with coverage" },

      -- Examples
      e = {
        name = "Examples",
        s = { "<cmd>!make run-simple<cr>", "Run simple example" },
        a = { "<cmd>!make run-animated<cr>", "Run animated example" },
        e = { "<cmd>!make run-editor<cr>", "Run editor example" },
      },

      -- Code Quality
      f = { "<cmd>!make fmt<cr>", "Format code" },
      l = { "<cmd>!make lint<cr>", "Lint code" },

      -- Documentation (MkDocs)
      D = {
        name = "Documentation",
        d = { "<cmd>!make docs-dev &<cr>", "Start MkDocs dev server" },
        b = { "<cmd>!make docs-build<cr>", "Build documentation" },
        c = { "<cmd>!make docs-clean<cr>", "Clean documentation" },
        o = { "<cmd>!make docs-open<cr>", "Open docs in browser" },
      },

      -- Demo Recording (asciinema)
      r = {
        name = "Recording",
        r = { "<cmd>terminal nix run .#demo-record<cr>", "Record demo" },
        p = { "<cmd>terminal nix run .#demo-play<cr>", "Play demo" },
      },

      -- Release
      R = { "<cmd>terminal nix run .#release<cr>", "Release new version" },

      -- Nix
      n = {
        name = "Nix",
        b = { "<cmd>!nix build<cr>", "Nix build" },
        r = { "<cmd>!nix run<cr>", "Nix run (demo)" },
        d = { "<cmd>!nix run .#docs-serve<cr>", "Nix docs serve" },
        t = { "<cmd>!nix run .#test<cr>", "Nix test" },
      },

      -- Version
      v = { "<cmd>!make version<cr>", "Show version info" },

      -- Dependencies
      u = { "<cmd>!make deps<cr>", "Update dependencies (go mod tidy)" },
      V = { "<cmd>!make vendor<cr>", "Vendor dependencies" },
    },
  }, { prefix = "<leader>" })
end

-- Autocommands for this project
local blockfont_group = vim.api.nvim_create_augroup("BlockfontProject", { clear = true })

-- Format Go files on save
vim.api.nvim_create_autocmd("BufWritePre", {
  group = blockfont_group,
  pattern = "*.go",
  callback = function()
    vim.lsp.buf.format({ async = false })
  end,
})

-- Run tests on test file save (optional, can be noisy)
-- Uncomment if you want automatic test runs
-- vim.api.nvim_create_autocmd("BufWritePost", {
--   group = blockfont_group,
--   pattern = "*_test.go",
--   callback = function()
--     vim.notify("Running tests...", vim.log.levels.INFO)
--     vim.fn.jobstart("go test -v ./...", {
--       on_exit = function(_, code)
--         if code == 0 then
--           vim.notify("Tests passed!", vim.log.levels.INFO)
--         else
--           vim.notify("Tests failed!", vim.log.levels.ERROR)
--         end
--       end,
--     })
--   end,
-- })

-- LSP configuration for gopls
vim.lsp.config.gopls = {
  settings = {
    gopls = {
      analyses = {
        unusedparams = true,
        shadow = true,
      },
      staticcheck = true,
      gofumpt = true,
    },
  },
}

-- Print confirmation
vim.notify("blockfont project configuration loaded", vim.log.levels.INFO)
