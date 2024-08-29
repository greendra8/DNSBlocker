document.addEventListener('DOMContentLoaded', () => {
    const userIdSpan = document.getElementById('userId');
    const optionsBtn = document.getElementById('optionsBtn');
  
    chrome.runtime.sendMessage({ action: 'getUserId' }, (response) => {
      userIdSpan.textContent = response.userId || 'Not set';
    });
  
    optionsBtn.addEventListener('click', () => {
      chrome.runtime.openOptionsPage();
    });
  });