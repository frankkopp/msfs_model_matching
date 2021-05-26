# msfs_model_matching
Simplified managing of model matching rules for VATSIM vPilot and MSFS

Main approach:

- read a list of all models in the sim with their default ICAO, their base model and the base livery (model name)
- read a list of variations of ICAO codes for planes which should share the base model
- read config path for Community folder
- read all folder with manifest.json and filter for folders where manifest.json has"content_type": "AIRCRAFT"
- AND has aircraft.cfg file
- LOOP over all remaining folders
    - read aircraft.cfg (SimObjects/Airplanes/*/aircraft.cfg)
    - read base_container
    - read title
    - read icao_airline
        - LOOP over all type code variations for this base model
        - write match rule entry with ICAO (CallSign Prefix), TypeCode, Livery Model Name
- 
    
