document.addEventListener('DOMContentLoaded', () => {
    const userIdInput = document.getElementById('userId');
    const saveBtn = document.getElementById('saveBtn');
    const status = document.getElementById('status');
    const refreshBtn = document.getElementById('refreshBtn');
    const blockedDomainsList = document.getElementById('blockedDomainsList');
    const countdownElement = document.getElementById('countdown');

    let countdownInterval;

    // Load user ID and blocked domains
    chrome.storage.local.get(['userId', 'blockedDomains', 'lastSyncTime'], (result) => {
        userIdInput.value = result.userId || '';
        displayBlockedDomains(result.blockedDomains || []);
        if (result.lastSyncTime) {
            startCountdown(result.lastSyncTime);
        }
    });

    saveBtn.addEventListener('click', () => {
        const newUserId = userIdInput.value;
        chrome.runtime.sendMessage({ action: 'setUserId', userId: newUserId }, () => {
            status.textContent = 'User ID saved';
            setTimeout(() => {
                status.textContent = '';
            }, 3000);
        });
    });

    refreshBtn.addEventListener('click', () => {
        chrome.runtime.sendMessage({ action: 'syncBlockedDomains' }, (response) => {
            if (response.success) {
                displayBlockedDomains(response.blockedDomains);
                status.textContent = 'Blocked domains refreshed';
                startCountdown(Date.now());
            } else {
                status.textContent = 'Error refreshing blocked domains';
            }
            setTimeout(() => {
                status.textContent = '';
            }, 3000);
        });
    });

    function displayBlockedDomains(domains) {
        blockedDomainsList.innerHTML = '';
        domains.forEach(domain => {
            const li = document.createElement('li');
            li.textContent = domain;
            blockedDomainsList.appendChild(li);
        });
    }

    function startCountdown(lastSyncTime) {
        clearInterval(countdownInterval);
        updateCountdown(lastSyncTime);
        countdownInterval = setInterval(() => updateCountdown(lastSyncTime), 1000);
    }

    function updateCountdown(lastSyncTime) {
        const now = Date.now();
        const nextSync = lastSyncTime + 5 * 60 * 1000; // 5 minutes in milliseconds
        const timeLeft = Math.max(0, nextSync - now);
        const minutes = Math.floor(timeLeft / 60000);
        const seconds = Math.floor((timeLeft % 60000) / 1000);
        countdownElement.textContent = `Next refresh in ${minutes}:${seconds.toString().padStart(2, '0')}`;
    }
});