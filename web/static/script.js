import { ethers, WebSocketProvider } from "https://cdnjs.cloudflare.com/ajax/libs/ethers/6.7.0/ethers.min.js";

const burnerKeyKey = "burnerKey";

// Retrieve or create a new burner key
function getOrCreateKey() {
    let key = localStorage.getItem(burnerKeyKey);
    if (!key) {
        key = ethers.Wallet.createRandom().privateKey;
        localStorage.setItem(burnerKeyKey, key);
    }
    return key;
}

// Update the status message displayed to the user
function setStatus(status) {
    document.getElementById("label-status").innerText = status;
}

// Create a new WebSockets provider for Ethereum
function newProvider(wsUrl) {
    return new WebSocketProvider(wsUrl);
}

// Retrieve or create a new burner wallet
function getOrCreateBurnerWallet(provider) {
    const key = getOrCreateKey();
    return new ethers.Wallet(key, provider);
}

// Log and display errors
async function catchError(error) {
    console.error('Error:', error);
    setStatus(`Error: ${error.message}`);
}

// Helper function to make POST requests with JSON data
async function postJson(url, data) {
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json', },
            body: JSON.stringify(data),
        });
        if (!response.ok) {
            throw new Error(await response.text());
        }
        return await response.json();
    } catch (error) {
        catchError(error);
    }
}

// Create a new lobby
async function newLobby(data) {
    console.log("Creating new lobby:", data);
    setStatus("Creating game lobby...");
    try {
        const resp = await postJson(`${location.origin}/actions/new-lobby`, data);
        setStatus("Lobby created.");
        return resp;
    } catch (error) {
        catchError(error);
    }
}

// Join an existing lobby
async function joinLobby(data) {
    console.log("Joining lobby:", data);
    setStatus("Joining game lobby...");
    try {
        const resp = await postJson(`${location.origin}/actions/join-lobby`, data);
        setStatus("Game joined.");
        return resp;
    } catch (error) {
        catchError(error);
    }
}

const factoryAbi = [
    "function createGame(string lobbyId, address[] _players) external returns (address)",
    "event GameCreated(address gameAddress, string lobbyId, address sender, address origin)"
];

// Listen for game creation event
async function listenForGameCreation(wallet, factoryAddress, providedLobbyId) {
    console.log("Listening for game creation:", { factoryAddress, providedLobbyId });
    const factory = new ethers.Contract(factoryAddress, factoryAbi, wallet);
    return new Promise((resolve) => {
        factory.on("GameCreated", (gameAddress, eventLobbyId) => {
            if (eventLobbyId === providedLobbyId) {
                console.log("Game created:", { gameAddress, eventLobbyId });
                resolve(gameAddress);
            }
        });
    });
}

// Create a game
async function createGame(wallet, factoryAddress, lobbyId, players) {
    console.log("Creating game:", { factoryAddress, lobbyId, players });
    setStatus("Creating game...");
    const factory = new ethers.Contract(factoryAddress, factoryAbi, wallet);
    try {
        const tx = await factory.createGame(lobbyId, players, { gasLimit: 2000000 });
        console.log("Transaction:", tx);
        setStatus("Waiting for transaction...");
        const receipt = await tx.wait();
        console.log("Receipt:", receipt);
        setStatus("Game created.");
        const gameCreatedEvent = receipt.logs.map(log => {
            try {
                return factory.interface.parseLog(log);
            } catch (error) {
                return null; // Ignore errors (logs that are not from this contract)
            }
        }).find(event => event && event.name === 'GameCreated');
        return gameCreatedEvent ? { gameAddress: gameCreatedEvent.args.gameAddress } : null;
    } catch (error) {
        catchError(error);
    }
}

// Request a drip of testnet Ether
async function requestDrip(address) {
    console.log("Requesting drip for address:", address);
    try {
        const resp = await postJson(`${location.origin}/actions/request-drip`, { address });
        if (resp && resp.ok) {
            console.log("Drip requested successfully:", resp);
        }
    } catch (error) {
        catchError(error);
    }
}

const tickMasterAbi = [
    "function totalGasAllocation() external view returns (uint256)",
    "function maxGasAllocation() external view returns (uint256)",
];

// Get the current utilization of the chain
async function getChainUtilization(provider, tickMasterAddress) {
    const tickMaster = new ethers.Contract(tickMasterAddress, tickMasterAbi, provider);
    const current = await tickMaster.totalGasAllocation();
    const max = await tickMaster.maxGasAllocation();
    const utilization = Number(current) / Number(max);
    return { current: Number(current), max: Number(max), utilization };
}

async function postHogInit() {
    !function (t, e) { var o, n, p, r; e.__SV || (window.posthog = e, e._i = [], e.init = function (i, s, a) { function g(t, e) { var o = e.split("."); 2 == o.length && (t = t[o[0]], e = o[1]), t[e] = function () { t.push([e].concat(Array.prototype.slice.call(arguments, 0))) } } (p = t.createElement("script")).type = "text/javascript", p.async = !0, p.src = s.api_host + "/static/array.js", (r = t.getElementsByTagName("script")[0]).parentNode.insertBefore(p, r); var u = e; for (void 0 !== a ? u = e[a] = [] : a = "posthog", u.people = u.people || [], u.toString = function (t) { var e = "posthog"; return "posthog" !== a && (e += "." + a), t || (e += " (stub)"), e }, u.people.toString = function () { return u.toString(1) + ".people (stub)" }, o = "capture identify alias people.set people.set_once set_config register register_once unregister opt_out_capturing has_opted_out_capturing opt_in_capturing reset isFeatureEnabled onFeatureFlags getFeatureFlag getFeatureFlagPayload reloadFeatureFlags group updateEarlyAccessFeatureEnrollment getEarlyAccessFeatures getActiveMatchingSurveys getSurveys onSessionId".split(" "), n = 0; n < o.length; n++)g(u, o[n]); e._i.push([i, s, a]) }, e.__SV = 1) }(document, window.posthog || []);
    const wallet = getOrCreateBurnerWallet()
    posthog.init('phc_LsLNl9nWm168DyAPmzYxLCG42C6dzWXGdSgQEmzfHx0', {
        api_host: "https://eu.posthog.com",
        persistence: 'memory',
        bootstrap: {
            distinctID: wallet.address,
        },
    });
}

export {
    setStatus,
    newProvider,
    getOrCreateBurnerWallet,
    newLobby,
    joinLobby,
    listenForGameCreation,
    createGame,
    requestDrip,
    getChainUtilization,
    postHogInit,
};
