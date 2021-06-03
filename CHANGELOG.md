## v0.3 (in progress)
- Windows UI
- DONE  
  - Make configuration savable
- ToDo:
  - Verify configuration at start
  - General: complete UI with menu, toolbar, etc.
  - Save should ask if existing file should be replaced (or create version copies)
  - Add custom data for skipping, ICAO changes, Livery changes, 
    - context click to edit livery custom data
  - Make configuration editable
  - Make cmd line run possible - no gui
  

  
## v0.2
- changed config format - separator : to ; (to be able to include drive letters)
- added configuration and handling for custom liveries (override or fixing of meta data)
- added cmd line option to show when a livery is custom (-showCustom)

## v0.1
- initial beta release
- features:
    - commandline args for search dir and configurations
    - creates VMR file base on liveries and config
