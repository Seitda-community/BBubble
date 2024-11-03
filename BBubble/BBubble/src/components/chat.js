// Enter 키로 메시지 전송
document.addEventListener("DOMContentLoaded", function() {
    const userInput = document.getElementById('user-input');
    if (userInput) {
        userInput.addEventListener('keypress', function(event) {
            if (event.key === "Enter") {
                sendMessage();
            }
        });
    }
});

// 메시지 전송 함수
function sendMessage() {
    const userInput = document.getElementById('user-input').value.trim();
    if (userInput === "") {
        alert("질문을 입력해 주세요.");
        return;
    }

    // 사용자 입력 메시지 추가
    addMessage("user", userInput);

    // GPT 응답 처리
    fetchGPTResponse(userInput);

    // 입력창 비우기
    document.getElementById('user-input').value = "";
}

// 메시지 추가 함수
function addMessage(sender, message) {
    const chatDisplay = document.getElementById('chat-display');

    const messageElement = document.createElement('p');
    messageElement.classList.add(sender);

    const spanElement = document.createElement('span');
    spanElement.innerText = message;

    messageElement.appendChild(spanElement);

    chatDisplay.appendChild(messageElement);
    chatDisplay.scrollTop = chatDisplay.scrollHeight;
}

// GPT 응답 처리 함수
function fetchGPTResponse(userInput) {
    addMessage("bot", "답변을 생성 중입니다...");
    
    setTimeout(() => {
        const fakeResponse = `${userInput}`;
        document.querySelector(".bot:last-child").innerText = fakeResponse;

        // 응답을 음성으로 전송
        sendVoice(fakeResponse);
    }, 1000);
}

// 음성 전송 함수 - Naver Clova API와 통신하여 음성으로 변환
function sendVoice(text) {
    const url = "https://naveropenapi.apigw.ntruss.com/voice/v1/tts";
    const options = {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
            "X-NCP-APIGW-API-KEY-ID": "YOUR_CLIENT_ID",
            "X-NCP-APIGW-API-KEY": "YOUR_CLIENT_SECRET"
        },
        body: `speaker=친근한&speed=0&text=${encodeURIComponent(text)}`
    };

    fetch(url, options)
        .then(response => response.blob())
        .then(blob => {
            const audioUrl = URL.createObjectURL(blob);
            const audio = new Audio(audioUrl);
            audio.play();
        })
        .catch(error => console.error("음성 전송 오류:", error));
}

// 마지막 봇 메시지를 업데이트하는 함수
function updateLastBotMessage(response) {
    const botMessages = document.querySelectorAll(".bot");
    if (botMessages.length > 0) {
        botMessages[botMessages.length - 1].innerText = response;
    }
}
