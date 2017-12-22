# AnemometerAlgorithm

This is intended to be the implementation of the Anemometer algorithm but at present the algorithm is implemented as part of the user interface. This is by default configured to output the results of the processing to standard output as well. The data structure looks:

Here is an example 
```json
{
  "Timestamp": 1513978578493813500,
  "Sensor": "579d5bca617c6479",
  "Vendor": "ucberkeley",
  "Algorithm": "1.0",
  "Tofs": [
    {
      "Src": 0,
      "Dst": 1,
      "Val": 555
    },
    {
      "Src": 0,
      "Dst": 2,
      "Val": 555
    },
    {
      "Src": 0,
      "Dst": 3,
      "Val": 555
    },
    {
      "Src": 1,
      "Dst": 0,
      "Val": 555
    },
    {
      "Src": 1,
      "Dst": 2,
      "Val": 555
    },
    {
      "Src": 1,
      "Dst": 3,
      "Val": 555
    },
    {
      "Src": 2,
      "Dst": 0,
      "Val": 555
    },
    {
      "Src": 2,
      "Dst": 1,
      "Val": 555
    },
    {
      "Src": 2,
      "Dst": 3,
      "Val": 555
    },
    {
      "Src": 3,
      "Dst": 0,
      "Val": 555
    },
    {
      "Src": 3,
      "Dst": 1,
      "Val": 555
    },
    {
      "Src": 3,
      "Dst": 2,
      "Val": 555
    }
  ],
  "RawInput": {
    "L7GHeaders": [
      {
        "Srcmac": "579d5bca617c6479",
        "Srcip": "fe80::559d:5bca:617c:6479",
        "Popid": "hk070",
        "Poptime": 74143819359,
        "Brtime": 1513978578493813500,
        "Rssi": 42,
        "Lqi": 255,
        "Payload": "CQtAtgDDND72AAAQHwK7AGsBEApeCA0MABucEswD7AAKA9//JwJl/zUBtP/pAGoArQA0AIAACgBVAOv/mP8JAZz/OQKu/pYDwvwcBA=="
      },
      {
        "Srcmac": "579d5bca617c6479",
        "Srcip": "fe80::559d:5bca:617c:6479",
        "Popid": "hk070",
        "Poptime": 74143986585,
        "Brtime": 1513978578661277000,
        "Rssi": 42,
        "Lqi": 255,
        "Payload": "CQxBtgHDMD72ACQQJAKwAG8BEgpeCA4JB/lkEkwEZQEuAxsBqAEVAbL/9ABp/7kAov91AOn/KAApAO3/F//lAff+JQFY/5oAtP+PAA=="
      },
      {
        "Srcmac": "579d5bca617c6479",
        "Srcip": "fe80::559d:5bca:617c:6479",
        "Popid": "hk070",
        "Poptime": 74144154489,
        "Brtime": 1513978578829213400,
        "Rssi": 42,
        "Lqi": 255,
        "Payload": "CQ1CtgLDCj7oACAQEwLIAG4BEgpeCAwIA72KEjj/mf93/7f/oP/e/7//BAC7AGQApABAAGMAHQARAP7/Wf80AGj/GgBp/w4AcP8BAA=="
      },
      {
        "Srcmac": "579d5bca617c6479",
        "Srcip": "fe80::559d:5bca:617c:6479",
        "Popid": "hk070",
        "Poptime": 74144321605,
        "Brtime": 1513978578996257000,
        "Rssi": 42,
        "Lqi": 255,
        "Payload": "CQpDtgPD/D3+AFAQJQK1AF0BEApeCAEHCX5jEiQA/gD3/1QCDwDqAzEASQX+/pwBtP82ATsAngB8AP3/GwCE/x0Abv8bAFz/EQBU/w=="
      }
    ],
    "ChirpHeaders": [
      {
        "Type": 11,
        "Seqno": 46656,
        "Build": 195,
        "CalPulse": 160,
        "CalRes": [
          4764,
          -1,
          -1,
          -1
        ],
        "Primary": 0,
        "MaxIndex": [
          -1,
          13,
          12,
          0
        ],
        "IValues": [
          null,
          [
            236,
            -33,
            -155,
            -76
          ],
          [
            106,
            52,
            10,
            -21
          ],
          [
            265,
            569,
            918,
            1052
          ]
        ],
        "QValues": [
          null,
          [
            972,
            778,
            551,
            309
          ],
          [
            233,
            173,
            128,
            85
          ],
          [
            -104,
            -100,
            -338,
            -830
          ]
        ],
        "Accelerometer": [
          -112.24,
          60.024,
          999.424
        ],
        "Magnetometer": [
          54.300000000000004,
          18.7,
          36.300000000000004
        ],
        "Temperature": 25.76,
        "Humidity": 21.42
      },
      {
        "Type": 12,
        "Seqno": 46657,
        "Build": 195,
        "CalPulse": 160,
        "CalRes": [
          -1,
          4708,
          -1,
          -1
        ],
        "Primary": 1,
        "MaxIndex": [
          14,
          -1,
          9,
          7
        ],
        "IValues": [
          [
            357,
            283,
            277,
            244
          ],
          null,
          [
            185,
            117,
            40,
            -19
          ],
          [
            485,
            293,
            154,
            143
          ]
        ],
        "QValues": [
          [
            1100,
            814,
            424,
            -78
          ],
          null,
          [
            -151,
            -94,
            -23,
            41
          ],
          [
            -233,
            -265,
            -168,
            -76
          ]
        ],
        "Accelerometer": [
          -113.216,
          60.024,
          1008.208
        ],
        "Magnetometer": [
          54.800000000000004,
          17.6,
          36.7
        ],
        "Temperature": 25.78,
        "Humidity": 21.42
      },
      {
        "Type": 13,
        "Seqno": 46658,
        "Build": 195,
        "CalPulse": 160,
        "CalRes": [
          -1,
          -1,
          4746,
          -1
        ],
        "Primary": 2,
        "MaxIndex": [
          12,
          8,
          -1,
          3
        ],
        "IValues": [
          [
            -103,
            -73,
            -34,
            4
          ],
          [
            100,
            64,
            29,
            -2
          ],
          null,
          [
            52,
            26,
            14,
            1
          ]
        ],
        "QValues": [
          [
            -200,
            -137,
            -96,
            -65
          ],
          [
            187,
            164,
            99,
            17
          ],
          null,
          [
            -167,
            -152,
            -151,
            -144
          ]
        ],
        "Accelerometer": [
          -122.488,
          56.608,
          1007.232
        ],
        "Magnetometer": [
          53.1,
          20,
          36.6
        ],
        "Temperature": 25.78,
        "Humidity": 21.42
      },
      {
        "Type": 10,
        "Seqno": 46659,
        "Build": 195,
        "CalPulse": 160,
        "CalRes": [
          -1,
          -1,
          -1,
          4707
        ],
        "Primary": 3,
        "MaxIndex": [
          1,
          7,
          9,
          -1
        ],
        "IValues": [
          [
            254,
            596,
            1002,
            1353
          ],
          [
            412,
            310,
            158,
            -3
          ],
          [
            -124,
            -146,
            -164,
            -172
          ],
          null
        ],
        "QValues": [
          [
            36,
            -9,
            15,
            49
          ],
          [
            -258,
            -76,
            59,
            124
          ],
          [
            27,
            29,
            27,
            17
          ],
          null
        ],
        "Accelerometer": [
          -125.904,
          61.976,
          1018.944
        ],
        "Magnetometer": [
          54.900000000000006,
          18.1,
          34.9
        ],
        "Temperature": 25.76,
        "Humidity": 21.42
      }
    ],
    "SetInfo": {
      "Site": "site0",
      "MAC": "579d5bca617c6479",
      "Build": 195,
      "Complete": true,
      "TimeOfFirst": "2017-12-22T13:36:18.49381362-08:00",
      "IsDuct": true
    }
  },
  "Temperatures": [
    {
      "Src": 0,
      "Dst": 1,
      "Val": 25.5
    },
    {
      "Src": 0,
      "Dst": 2,
      "Val": 25.5
    },
    {
      "Src": 0,
      "Dst": 3,
      "Val": 25.5
    },
    {
      "Src": 1,
      "Dst": 0,
      "Val": 25.5
    },
    {
      "Src": 1,
      "Dst": 2,
      "Val": 25.5
    },
    {
      "Src": 1,
      "Dst": 3,
      "Val": 25.5
    },
    {
      "Src": 2,
      "Dst": 0,
      "Val": 25.5
    },
    {
      "Src": 2,
      "Dst": 1,
      "Val": 25.5
    },
    {
      "Src": 2,
      "Dst": 3,
      "Val": 25.5
    },
    {
      "Src": 3,
      "Dst": 0,
      "Val": 25.5
    },
    {
      "Src": 3,
      "Dst": 1,
      "Val": 25.5
    },
    {
      "Src": 3,
      "Dst": 2,
      "Val": 25.5
    }
  ],
  "Velocities": {
    "X": 1,
    "Y": 2,
    "Z": 3,
    "Mag": 3.7416573867739413,
    "Phi": 55,
    "Theta": 66,
    "Valid": true
  },
  "Extradata": [
    "the algorithm has not been filled in yet"
  ],
  "Uncorrectable": 1,
  "Correctable": 0,
  "Total": 17
}
```
