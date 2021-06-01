# MSFS VATSIM vPilot MatchMaker
[![Go](https://github.com/frankkopp/msfs_model_matching/actions/workflows/go.yml/badge.svg)](https://github.com/frankkopp/msfs_model_matching/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/frankkopp/msfs_model_matching)](https://goreportcard.com/report/github.com/frankkopp/msfs_model_matching)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/frankkopp/WorkerPool/blob/master/LICENSE)

Automatically generates a vPilot match making file (VMR) for all liveries in the given folder.
It allows to configure the base sim plane models to be used, the default liveries to be used 
and also a list of alternative ICAO codes for airline which should use the same livery.

## Installation
- Download the latest release or the most current master build (https://github.com/frankkopp/msfs_model_matching/releases/tag/latest_build)
- Extract this to a folder (I use the vPilot files folder in my Documents folder)

## Configuration
There are currently four configuration files required.
(these might be consolidated to less in the future)

### defaultTypes.txt

This maps the sims base container each livery has to reference to and their respective 
default livery.
As the sim will never have all real life plane models this is a way to still 
show other planes with a default livery if there is otherwise no matching livery available. 
It is possible to configure several default liveries and vPilot will choose randomly.
````
Asobo_A320_NEO;     Airbus A320 Neo Asobo
Asobo_B747_8i;      Boeing 747-8i Asobo
Asobo_B787_10;      Boeing 787-10 Asobo
Aerosoft_CRJ_700;   CRJ550ER Privat D-ALKI
Asobo_CJ4;          Cessna CJ4 Citation Asobo
Asobo_Longitude;    Cessna Citation Longitude Asobo
Asobo_TBM930;       TBM 930 Asobo;TBM 930 Asobo Air Traffic 00;TBM 930 Asobo Air Traffic 01;TBM 930 Asobo Air Traffic 02;
````

### typeVariations.txt
This maps the sims base container each livery has to reference and all plane type codes 
which should use the same liveries.
As the sim will never have all real life plane models this is a way to still match
planes used in VATSIM by other pilots. 
````
# Narrow/Medium Jet 2 Engines
Asobo_A320_NEO;     A19N;A20N;A21N;A318;A319;A320;A321;B732;B733;B734;B735;B736;B737;B738;B739;B73X;B37M;B38M;B39M

# Wide/Heavy Jet 2 Engines
Asobo_B787_10;      B78X;B788;B789;B762;B763;B764;B772;B773;B778;B779;B77L;B77W;A306;A30B;A310;A332;A333;A337;A338;A339

# Wide/Heavy Jet 4 Engines
Asobo_B747_8i;      B741;B742;B743;B744;B748;B74F;A380;A388

# Small/Light Jet 2 Tail Engines
Aerosoft_CRJ_700;   CRJ7;CRJX;CRJ5;CRJ9

# Business Jet 2 Tail Engines
Asobo_CJ4;          C25C;C25B;C25A;C500;C501;C510;C525;C526
Asobo_Longitude;    C700;C750

# Turbo Prop
Asobo_TBM930;       TBM9
````

### icaoVariations.txt
ICAO airline codes typically have variations which still belong to the 
same airline. Mapping them as equivalent here allows rules to be generated
which will map these ICAO codes to the same livery as the ICAO code which
was part of the livery definition. Therefore, it is not necessary to have a 
livery for each of the variations. 
E.g. British Airways has an ICAO code of "BAW". But, sometimes pilots also use
BA, SHT or CFE (etc.) - all these planes will use the BAW livery if so configured. 
````
DLH;LHA;CLH
BAW;BA;SHT;CFE
EZY;EJU;EZS
TUI;TOM;TFL;THOM
EWG;EWE
RYR;RUK
DHK;BCS;DHL;DAE
````

### customData.txt
This file can be used to configure custom data for a livery (aircraft.cfg).
The required format is path;name;base;icao
For skipped files this can easily be copied from the output of the program and 
pasted into the file. Remove SKIPPED at the beginning and enter the name, base and icao manually.
Also, liveries can be explicitly skipped by using "skip" as name. 

````
# Format: path;name;base;icao
# overwrite with custom
D:\Games\MSFS2020\Community\Asobo_B747_8i_FedEx\SimObjects\Airplanes\Asobo_B747_8i-FedEx\aircraft.cfg;Boeing 747-8i FedEx;Asobo_B747_8i;FDX
# fix missing data
D:\Games\MSFS2020\Community\Aerosoft_CRJ_ACJazz\SimObjects\AirPlanes\Aerosoft_CRJ_700_JAZZ\aircraft.cfg;CRJ700ER Jazz Lowercase;Aerosoft_CRJ_700;JZA
D:\Games\MSFS2020\Community\Asobo_A320_NEO_AirCaraibes\SimObjects\Airplanes\Asobo_A320_NEO-AirCaraibes\aircraft.cfg;Airbus A320 Neo Air Caraibes;Asobo_A320_NEO;FWI
# skip these
D:\Games\MSFS2020\Community\aircraft-longitudeFDEfix\SimObjects\AirPlanes\Asobo_Longitude\aircraft.cfg;skip;;
D:\Games\MSFS2020\Community\aircraft-tbm930x\SimObjects\Airplanes\Asobo_TBM930\aircraft.cfg;skip;;
````

## Usage
matchmaker.exe is used from the command line. So open a terminal and go to 
the folder where the exe file is.
Make sure you have a folder config/ with the config files or use command line arguments 
to define where these files are. 

Run:
````
matchmaker.exe -dir <path to your msfs community folder>
````

The resulting file will be stored in the same folder as the exe with the filename (or use the -outputFile parameter):
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
  -fixLiveriesFile string
        path and filename to fix liveries config file (default "..\\config\\fixLiveries.txt")
  -icaoVariationsFile string
        path and filename to icao variations config file (default "..\\config\\icaoVariations.txt")
  -outputFile string
        path and filename to output file (default ".\\MatchMakingRules.vmr")
  -showCustom
        shows liveries which are custom by configuration
  -typeVariationsFile string
        path and filename to type variations config file (default "..\\config\\typeVariations.txt")
  -version
        prints version and exits

````

## How it works:

When started MatchMaker.exe searches recursively for aircraft.cfg files in the given folder
- aircraft.cfg must contain these three data points otherwise it will be skipped:
    - base_container
    - icao_airline
    - title

For each found aircraft.cfg (aka Livery) rules will be created based on the ICAO code of the livery if the
base_container is a model which is configured and has at least one mapping to a plane type code.
If there are several liveries for an ICAO code, and a type code the liveries will be concatenated, so vPilot can
choose randomly which to use.
If metadata is missing, or a custom values are required the aircraft.cfg file can be added to the customData.txt 
config file. 

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
