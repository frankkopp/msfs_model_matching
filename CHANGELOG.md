## v1.1 
- Planes are recognized as well, not only pure liveries
- Prompt user to save configuration and/or rules on window close
- Create backup when saving rules file

## v1.0
- Windows UI
- Make configuration savable
- Make configuration editable (simple version)
- Context click to edit livery custom data (one item)
- Add custom data for skipping, ICAO changes, Livery changes, 
- Context click to edit livery custom data (multiple items)
- CleanUp and restructure 
- Save should ask if existing file should be replaced (or create version copies)
- configuration dirty logic for status bar
- Make cmd line run possible - no gui
- remember window position
    
## v0.2
- changed config format - separator : to ; (to be able to include drive letters)
- added configuration and handling for custom liveries (override or fixing of metadata)
- added cmd line option to show when a livery is custom (-showCustom)

## v0.1
- initial beta release
- features:
    - commandline args for search dir and configurations
    - creates VMR file base on liveries and config
