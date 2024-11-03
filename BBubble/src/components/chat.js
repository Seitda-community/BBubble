// 챗봇 메시지 전송 함수
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

    // p 요소 생성
    const messageElement = document.createElement('p');
    messageElement.classList.add(sender);

    // span 요소 생성 및 메시지 추가
    const spanElement = document.createElement('span');
    spanElement.innerText = message;

    // p 안에 span 추가
    messageElement.appendChild(spanElement);

    // chat-display에 메시지 추가
    chatDisplay.appendChild(messageElement);
    chatDisplay.scrollTop = chatDisplay.scrollHeight; // 자동 스크롤
}

// GPT 응답 처리 함수
function fetchGPTResponse(userInput) {
	addMessage("bot", "답변을 생성 중입니다...");
	
	// 여기에서 GPT API 호출 로직을 추가할 수 있습니다.
	setTimeout(() => {
		const fakeResponse = `${userInput}`; // 테스트용 응답
		document.querySelector(".bot:last-child").innerText = fakeResponse;
	}, 2000);
}

// 마지막 봇 메시지를 업데이트하는 함수
function updateLastBotMessage(response) {
    const botMessages = document.querySelectorAll(".bot");
    if (botMessages.length > 0) {
        botMessages[botMessages.length - 1].innerText = response;
    }
}