# MSFS VATSIM vPilot MatchMaker
Automatically generates a vPilot match making file (VMR) for all liveries in the given folder.
It allows to configure the base sim plane models to be used, the default liveries to be used 
and also a list of alternative ICAO codes for airline which should use the same livery.

## Installation
- Download the latest release or the most current master build (https://github.com/frankkopp/msfs_model_matching/releases/tag/latest_build)
- Extract this to a folder (I use the vPilot files folder in my Documents folder)

## Configuration
There are currently three configuration files required.
(these might be consolidated to 2 in the future)

### defaultTypes.txt

This maps the sims base container the livery has to point to and their respective 
default livery.
As the sim will never have all real life plane models this is a way to still 
show other planes in nice default livery if there is no livery available. 
````
Asobo_A320_NEO:     Airbus A320 Neo Asobo
Asobo_B747_8i:      Boeing 747-8i Asobo
Asobo_B787_10:      Boeing 787-10 Asobo
Aerosoft_CRJ_700:   CRJ550ER Privat D-ALKI
Asobo_CJ4:          Cessna CJ4 Citation Asobo
Asobo_Longitude:    Cessna Citation Longitude Asobo
Asobo_TBM930:       TBM 930 Asobo
````

### typeVariations.txt
This maps the base container the livery has to point to and all plane type codes 
which should use this livery.
As the sim will never have all real life plane models this is a way to still match
planes used in VATSIM by other pilots. 
````
Asobo_A320_NEO:     A19N:A20N:A21N:A318:A319:A320:A321:B732:B733:B734:B735:B736:B737:B738:B739:B73X
Asobo_B747_8i:      B741:B742:B743:B744:B748:B74F
Asobo_B787_10:      B78X:B788:B789
Asobo_CJ4:          C25C:C25B:C25A:C500:C501:C510:C525:C526
Asobo_Longitude:    C700:C750
Aerosoft_CRJ_700:   CRJ7:CRJX:CRJ5:CRJ9
Asobo_TBM930:       TBM9
````

### icaoVariations.txt
ICAO airline codes typically have variations which still belong to the 
same airline. Mapping them as equivalent here allows rule to be generated.
which will map these ICAO codes to the same livery then the ICAO code which
was part of the livery definition. Therefore, it is not necessary to have an 
own livery for each of the variations. 
E.g. British Airways has an ICAO code of "BAW". But sometimes pilots also use
BA, SHT or CFE (etc.) - all these planes will use the BAW livery if so configured. 
````
DLH:LHA:CLH
BAW:BA:SHT:CFE
EZY:EJU:EZS
TUI:TOM:TFL:THOM
EWG:EWE
RYR:RUK
DHK:BCS:DHL:DAE
````

## Usage
matchmaker.exe is used from the command line. So open a terminal and go to 
the folder where the exe file is.
Make sure you have a folder config/ with the three config files or use command line arguments 
to define where these files are. 

Run:
````
matchmaker.exe -dir <path to your msfs community folder>
````

The resulting file will be stored in the same folder as the exe with the filename:
````
MatchMakingRules.vmr
````

Use this file within vPilot as Model matching rules file. 

### Command line options

````
Usage of matchmaker.exe:
-defaultTypesFile string 
    path and filename to default types config file (default "..\\config\\defaultTypes.txt")
-dir string 
    path where aircraft.cfg are searched recursively (default ".")
-icaoVariationsFile string
    path and filename to icao variations config file (default "..\\config\\icaoVariations.txt")
-outPutFile string
    path and filename to output file (default ".\\MatchMakingRules.vmr")
-typeVariationsFile string
    path and filename to type variations config file (default "..\\config\\typeVariations.txt")
-version
    prints version and exits
````


## How it works:

Searches recursively for aircraft.cfg file
- aircraft.cfg must contain these three data points otherwise it will be ignored:
    - base_container
    - icao_airline
    - title

For each found aircraft.cfg (aka Livery) rules will be created based on the ICAO code of the livery if the
base_container is a model which is configured and has at least one mapping to a plane type code.
If there are several liveries for an ICAO code, and a type code the liveries will be concatenated, so vPilot can
choose randomly which to use.

Example:
- Livery for A320Neo for Lufthansa (DLH)
- Base container for A320Neo is: Asobo_A320_NEO
- TypeCodes for this base container are: A19N:A20N:A21N:A318:A319:A320:A321:B732:B733:B734:B735:B736:B737:B738:B739:B73X
- Alternative ICAO airline codes for DLH are: DLH:LHA:CLH

Generated rules:
````
<!-- AIRLINE ICAO: CLH -->
<!-- BASE TYPE: Asobo_A320_NEO -->
<ModelMatchRule CallsignPrefix="CLH" TypeCode="A19N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="A20N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="A21N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="A318" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="A319" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="A320" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="A321" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B732" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B733" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B734" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B735" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B736" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B737" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B738" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B739" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="CLH" TypeCode="B73X" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />

<!-- AIRLINE ICAO: DLH -->
<!-- BASE TYPE: Asobo_A320_NEO -->
<ModelMatchRule CallsignPrefix="DLH" TypeCode="A19N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="A20N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="A21N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="A318" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="A319" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="A320" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="A321" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B732" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B733" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B734" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B735" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B736" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B737" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B738" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B739" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="DLH" TypeCode="B73X" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />

<!-- BASE TYPE: Asobo_A320_NEO -->
<ModelMatchRule CallsignPrefix="LHA" TypeCode="A19N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="A20N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="A21N" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="A318" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="A319" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="A320" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="A321" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B732" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B733" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B734" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B735" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B736" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B737" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B738" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B739" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
<ModelMatchRule CallsignPrefix="LHA" TypeCode="B73X" ModelName="Airbus A320 Neo Lufthansa AI OperatorLand468//Airbus A320 Neo Lufthansa Gummersbach OperatorLand468//Airbus A320 Neo Lufthansa Neubrandenburg OperatorLand468" />
````
