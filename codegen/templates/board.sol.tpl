// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;
{{ range .Imports }}
import "{{ . }}";
{{- end }}

uint8 constant WIDTH = {{ $.Width }};
uint8 constant HEIGHT = {{ $.Height }};

library BoardLib {
    function initCore(ICore proxy, uint16 w, uint16 h) internal {
        ActionData_Initialize memory initializeData;
        initializeData.width = w;
        initializeData.height = h;
        proxy.initialize(initializeData);
    }

    function addBuilding(ICore proxy, uint8 playerId, uint8 buildingTypeId, uint16 x, uint16 y) internal {
        ActionData_PlaceBuilding memory placeBuildingData;
        placeBuildingData.playerId = playerId;
        placeBuildingData.buildingType = buildingTypeId;
        placeBuildingData.x = x;
        placeBuildingData.y = y;
        proxy.placeBuilding(placeBuildingData);
    }

    function addUnit(ICore proxy, uint8 playerId, uint8 unitTypeId, uint16 x, uint16 y) internal {
        ActionData_CreateUnit memory createUnitData;
        createUnitData.playerId = playerId;
        createUnitData.unitType = unitTypeId;
        createUnitData.x = x;
        createUnitData.y = y;
        proxy.createUnit(createUnitData);
    }
    
    function assignUnit(ICore proxy, uint8 playerId, uint8 unitId, uint64 command) internal {
        ActionData_AssignUnit memory assignUnitData;
        assignUnitData.playerId = playerId;
        assignUnitData.unitId = unitId;
        assignUnitData.command = command;
        proxy.assignUnit(assignUnitData);
    }

    function initPlayer(ICore proxy, uint8 playerId, uint8 unpurgeableUnitCount) internal {
        ActionData_AddPlayer memory addPlayerData;
        addPlayerData.unpurgeableUnitCount = unpurgeableUnitCount;

        {{ range $playerId, $player := $.Players -}}
        {{- if ne $playerId 0 -}}
        if (playerId == {{ $playerId }}) {
            addPlayerData.spawnAreaX = {{ $player.SpawnArea.Min.X }};
            addPlayerData.spawnAreaY = {{ $player.SpawnArea.Min.Y }};
            addPlayerData.spawnAreaWidth = {{ $player.SpawnArea.Dx }};
            addPlayerData.spawnAreaHeight = {{ $player.SpawnArea.Dy }};
            addPlayerData.workerPortX = {{ $player.WorkerPort.X }};
            addPlayerData.workerPortY = {{ $player.WorkerPort.Y }};
            proxy.addPlayer(addPlayerData);
            {{- range $i, $b := $player.Buildings }}
            addBuilding(proxy, {{ $playerId }}, {{ $b.PrototypeId }}, {{ $b.Position.X }}, {{ $b.Position.Y }});
            {{- end }}
            {{- range $i, $u := $player.Units }}
            addUnit(proxy, {{ $playerId }}, {{ $u.PrototypeId }}, {{ $u.Position.X }}, {{ $u.Position.Y }});
            {{- if $u.Command }}
            assignUnit(proxy, {{ $playerId }}, {{ add $i 1 }}, {{ $u.Command }});
            {{- end -}}
            {{- end }}
        } else {{ end }}{{ end }}{
            revert();
        }
    }

    function initEnvironment(ICore proxy) internal {
        {{- if index $.Players 0 -}}
        {{- range $i, $b := (index $.Players 0).Buildings }}
        addBuilding(proxy, 0, {{ $b.PrototypeId }}, {{ $b.Position.X }}, {{ $b.Position.Y }});
        {{- end -}}
        {{- end }}
    }

    function initialize(ICore proxy) internal {
        initCore(proxy, WIDTH, HEIGHT);
        initEnvironment(proxy);
    }
}
