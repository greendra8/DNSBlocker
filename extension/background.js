let blockedDomains = [];
let userId = '';
const API_BASE_URL = 'https://localhost:443'; // Update this to your CoreDNS server address

// Fetch blocked domains from server
async function syncBlockedDomains() {
  if (!userId) return;
  try {
    const response = await fetch(`${API_BASE_URL}/rules/${userId}`);
    const rules = await response.json();
    blockedDomains = rules.map(rule => rule.domain);
    const lastSyncTime = Date.now();
    chrome.storage.local.set({ blockedDomains, lastSyncTime });
    console.log('Synced blocked domains:', blockedDomains);
    return blockedDomains;
  } catch (error) {
    console.error('Error syncing blocked domains:', error);
    throw error;
  }
}

// Set up periodic sync
chrome.alarms.create('syncBlockedDomains', { periodInMinutes: 5 });
chrome.alarms.onAlarm.addListener(syncBlockedDomains);

// Intercept requests
chrome.webRequest.onBeforeRequest.addListener(
  (details) => {
    const url = new URL(details.url);
    if (blockedDomains.some(domain => url.hostname.endsWith(domain))) {
      return { redirectUrl: chrome.runtime.getURL('blocked.html') };
    }
  },
  { urls: ['<all_urls>'] },
  ['blocking']
);

// Load user ID and initial blocked domains
chrome.storage.local.get(['userId', 'blockedDomains'], (result) => {
  userId = result.userId;
  blockedDomains = result.blockedDomains || [];
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
      .then(() => {
        sendResponse({ success: true, blockedDomains });
      })
      .catch((error) => {
        console.error('Error syncing blocked domains:', error);
        sendResponse({ success: false });
      });
    return true; // Indicates that the response is sent asynchronously
  }
});