/*
 * MatchMaker - create model matching files for VATSIM vPilot
 *
 *  MIT License
 *
 *  Copyright (c) 2021 Frank Kopp
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 *
 */

package config

var defaultIni = `
[paths]
liveryDir = .
outputFile = .\MatchMakingRulesUI.vmr

[defaultTypes]
Asobo_A320_NEO = Airbus A320 Neo Asobo
Asobo_B747_8i = Boeing 747-8i Asobo
Asobo_B787_10 = Boeing 787-10 Asobo
# Aerosoft_CRJ_700 = CRJ550ER Privat D-ALKI
Asobo_CJ4 = Cessna CJ4 Citation Asobo
Asobo_Longitude = Cessna Citation Longitude Asobo
Asobo_TBM930 = TBM 930 Asobo,TBM 930 Asobo Air Traffic 00,TBM 930 Asobo Air Traffic 01,TBM 930 Asobo Air Traffic 02

[typeVariations]
# Narrow/Medium Jet 2 Engines
Asobo_A320_NEO = A19N,A20N,A21N,A318,A319,A320,A321,B732,B733,B734,B735,B736,B737,B738,B739,B73X,B37M,B38M,B39M
# Wide/Heavy Jet 2 Engines
Asobo_B787_10 = B78X,B788,B789,B762,B763,B764,B772,B773,B778,B779,B77L,B77W,A306,A30B,A310,A332,A333,A337,A338,A339
# Wide/Heavy Jet 4 Engines
Asobo_B747_8i = B741,B742,B743,B744,B748,B74F,A380,A388
# Small/Light Jet 2 Tail Engines
# Aerosoft_CRJ_700 = CRJ7,CRJX,CRJ5,CRJ9
# Business Jet 2 Tail Engines
Asobo_CJ4 = C25C,C25B,C25A,C500,C501,C510,C525,C526
Asobo_Longitude = C700,C750
# Turbo Prop
Asobo_TBM930 = TBM9

[icaoVariations]
Lufthansa = DLH,LHA,CLH
BritishAirways = BAW,BA,SHT,CFE
EasyJet = EZY,EJU,EZS
TUI = TUI,TOM,TFL,THOM
Eurowings = EWG,EWE
Ryanair = RYR,RUK
DHL = BCS,DHL,DAE,DHK
Fedex = FDX,FEDEX
Luxair = LUX,LGL
WizzAir = WZZ,WUK
VirginAtlantic = VIR,VOZ

# this section is automatically managed by the UI - edit with care
[customData]
Do not delete this line due to a bug in the ini library,false,,
`
