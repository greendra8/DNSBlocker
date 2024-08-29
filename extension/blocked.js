function getQueryParam(param) {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get(param);
}

function tryAgain() {
  const blockedDomain = getQueryParam('domain');
  if (blockedDomain) {
    window.location.href = `https://${decodeURIComponent(blockedDomain)}`;
  } else {
    window.close();
  }
}

function typeWriter(text, elementId, speed = 50) {
  const element = document.getElementById(elementId);
  element.innerHTML = '';
  let i = 0;
  function type() {
    if (i < text.length) {
      element.innerHTML += text.charAt(i);
      i++;
      setTimeout(type, speed);
    }
  }
  type();
}

const funnyMessages = [
  "Cosmic Productivity Shield: Engaged",
  "Temporal Distortion Detected: Redirecting to Success Continuum",
  "Quantum Focus Field: Stabilized",
  "Neural Pathways Optimized: Distraction Firewall Active",
  "Holographic Projection: Success Trajectory Locked",
  "Entering Hyper-Productive Space-Time Continuum",
  "Future Self Gratitude Protocol: Initiated",
  "Reality Distortion Field: Recalibrated for Peak Performance",
  "Laser Focus Module: Online",
  "Cosmic Knowledge Gateway: Temporarily Sealed"
];

document.addEventListener('DOMContentLoaded', function() {
  const blockedDomain = getQueryParam('domain');
  const blockMessage = document.getElementById('blockMessage');
  const funnyMessageElement = document.getElementById('funnyMessage');
  const tryAgainButton = document.getElementById('tryAgainButton');

  if (blockedDomain) {
    typeWriter(`Alert: ${decodeURIComponent(blockedDomain).toUpperCase()} lies beyond the cosmic productivity barrier.`, 'blockMessage');
  } else {
    typeWriter("Anomaly Detected: Unauthorized access to restricted cosmic sector.", 'blockMessage');
  }

  const randomIndex = Math.floor(Math.random() * funnyMessages.length);
  funnyMessageElement.textContent = funnyMessages[randomIndex];

  tryAgainButton.addEventListener('click', tryAgain);

  // Easter egg: Konami Code
  const konamiCode = ['ArrowUp', 'ArrowUp', 'ArrowDown', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'ArrowLeft', 'ArrowRight', 'b', 'a'];
  let konamiIndex = 0;

  document.addEventListener('keydown', (e) => {
    if (e.key === konamiCode[konamiIndex]) {
      konamiIndex++;
      if (konamiIndex === konamiCode.length) {
        document.body.style.animation = 'cosmicShift 3s linear';
        setTimeout(() => {
          document.body.style.animation = 'none';
        }, 3000);
        konamiIndex = 0;
      }
    } else {
      konamiIndex = 0;
    }
  });
});

// Add this to your CSS in the HTML file
document.head.insertAdjacentHTML('beforeend', `
  <style>
    @keyframes cosmicShift {
      0% { filter: hue-rotate(0deg) brightness(1); }
      50% { filter: hue-rotate(180deg) brightness(1.5); }
      100% { filter: hue-rotate(360deg) brightness(1); }
    }
  </style>
`);