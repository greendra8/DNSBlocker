let blockedDomains = [];
let userId = '';
const API_BASE_URL = 'https://localhost:443'; // Update this to your CoreDNS server address
let lastSyncTime = 0;

// Fetch blocked domains from server
async function syncBlockedDomains() {
  if (!userId) {
    console.log('No user ID set, skipping sync');
    return;
  }
  try {
    const response = await fetch(`${API_BASE_URL}/rules/${userId}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const rules = await response.json();
    console.log('Raw response:', rules);

    if (Array.isArray(rules)) {
      blockedDomains = rules.map(rule => rule.domain);
    } else {
      blockedDomains = []; // Set to empty array if no rules found
    }
    
    lastSyncTime = Date.now();
    chrome.storage.local.set({ blockedDomains, lastSyncTime });
    console.log('Synced blocked domains:', blockedDomains);
    return blockedDomains;
  } catch (error) {
    console.error('Error syncing blocked domains:', error);
    console.error('Error details:', error.message);
    throw error;
  }
}

// Set up periodic sync
chrome.alarms.create('syncBlockedDomains', { periodInMinutes: 1 });  // Changed from 5 to 1
chrome.alarms.onAlarm.addListener(syncBlockedDomains);

// Intercept requests
chrome.webRequest.onBeforeRequest.addListener(
  (details) => {
    const url = new URL(details.url);
    if (blockedDomains.some(domain => url.hostname.endsWith(domain))) {
      const encodedDomain = encodeURIComponent(url.hostname);
      return { redirectUrl: chrome.runtime.getURL(`blocked.html?domain=${encodedDomain}`) };
    }
  },
  { urls: ['<all_urls>'] },
  ['blocking']
);

// Load user ID and initial blocked domains
chrome.storage.local.get(['userId', 'blockedDomains', 'lastSyncTime'], (result) => {
  userId = result.userId;
  blockedDomains = result.blockedDomains || [];
  lastSyncTime = result.lastSyncTime || 0;
  if (userId) syncBlockedDomains();
});

// Listen for messages from popup
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'getUserId') {
    sendResponse({ userId });
  } else if (request.action === 'setUserId') {
    userId = request.userId;
    chrome.storage.local.set({ userId });
    syncBlockedDomains();
  } else if (request.action === 'syncBlockedDomains') {
    syncBlockedDomains()
      .then((domains) => {
        sendResponse({ 
          success: true, 
          blockedDomains: domains,
          message: domains.length > 0 ? `Found ${domains.length} blocked domain(s)` : 'No blocked domains found'
        });
      })
      .catch((error) => {
        console.error('Error syncing blocked domains:', error);
        sendResponse({ success: false, message: 'Error syncing blocked domains' });
      });
    return true; // Indicates that the response is sent asynchronously
  } else if (request.action === 'getLastSyncTime') {
    sendResponse({ lastSyncTime });
  }
});

const backgrounds = [
  chrome.runtime.getURL('bg/moon.png'),
  chrome.runtime.getURL('bg/city.png'),
];

function setRandomBackground() {
  const randomBg = backgrounds[Math.floor(Math.random() * backgrounds.length)];
  const backgroundElement = document.querySelector('.background');
  backgroundElement.style.backgroundImage = `url('${randomBg}')`;
  console.log('Background set to:', randomBg); // For debugging
}

setRandomBackground();
window.addEventListener('load', setRandomBackground);