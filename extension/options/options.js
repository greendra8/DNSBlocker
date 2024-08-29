document.addEventListener('DOMContentLoaded', () => {
    const userIdInput = document.getElementById('userId');
    const saveBtn = document.getElementById('saveBtn');
    const status = document.getElementById('status');
    const refreshBtn = document.getElementById('refreshBtn');
    const blockedDomainsList = document.getElementById('blockedDomainsList');
    const countdownElement = document.getElementById('countdown');

    let countdownInterval;

    function loadData() {
        chrome.storage.local.get(['userId', 'blockedDomains'], (result) => {
            userIdInput.value = result.userId || '';
            displayBlockedDomains(result.blockedDomains || []);
        });
        startCountdown();
    }

    loadData();

    saveBtn.addEventListener('click', () => {
        const newUserId = userIdInput.value;
        chrome.runtime.sendMessage({ action: 'setUserId', userId: newUserId }, () => {
            status.textContent = 'User ID saved';
            setTimeout(() => {
                status.textContent = '';
            }, 3000);
        });
    });

    refreshBtn.addEventListener('click', refreshBlockedDomains);

    function refreshBlockedDomains() {
        chrome.runtime.sendMessage({ action: 'syncBlockedDomains' }, (response) => {
            if (response.success) {
                displayBlockedDomains(response.blockedDomains);
                status.textContent = response.message;
                startCountdown();
            } else {
                status.textContent = response.message;
            }
            setTimeout(() => {
                status.textContent = '';
            }, 3000);
        });
    }

    function displayBlockedDomains(domains) {
        blockedDomainsList.innerHTML = '';
        domains.forEach(domain => {
            const li = document.createElement('li');
            li.textContent = domain;
            blockedDomainsList.appendChild(li);
        });
    }

    function startCountdown() {
        clearInterval(countdownInterval);
        updateCountdown();
        countdownInterval = setInterval(updateCountdown, 1000);
    }

    function updateCountdown() {
        chrome.runtime.sendMessage({ action: 'getLastSyncTime' }, (response) => {
            const now = Date.now();
            const nextSync = response.lastSyncTime + 60 * 1000; // 1 minute in milliseconds
            const timeLeft = Math.max(0, nextSync - now);
            const seconds = Math.floor(timeLeft / 1000);
            countdownElement.textContent = `Next refresh in ${seconds} seconds`;
            
            if (seconds === 0) {
                refreshBlockedDomains();
            }
        });
    }
});