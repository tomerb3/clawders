-- Pull in the wezterm API
local wezterm = require 'wezterm'

-- This will hold the configuration.
local config = wezterm.config_builder()

-- For example, changing the initial geometry for new windows:
config.initial_cols = 120
config.initial_rows = 28

-- Font settings
config.font_size = 16
config.line_height = 1.2
-- config.font = wezterm.font_with_fallback { "3270 Nerd Font Mono", "Cascadia Code" }
config.font = wezterm.font_with_fallback { "AudioLink Mono", "Cascadia Code" }

-- Color scheme
config.color_scheme = 'tokyonight_night'

-- Startup
config.default_domain = 'WSL:Ubuntu'

-- Exit behavior - close without asking
config.exit_behavior = 'Close'
config.window_close_confirmation = 'NeverPrompt'

-- Cursor
config.default_cursor_style = "BlinkingBlock"
config.colors = {
	cursor_bg = "#FFD700",
	cursor_border = "#FFD700",
}

-- Show tab bar always
config.enable_tab_bar = true
config.hide_tab_bar_if_only_one_tab = false


-- Disable bell sound
config.audible_bell = "Disabled"

-- Hyperlinks
config.hyperlink_rules = wezterm.default_hyperlink_rules()

-- Finally, return the configuration to wezterm:
return config
