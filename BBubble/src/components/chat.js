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

// 타자기 효과 함수
function typeWriterEffect(htmlText, element, speed = 50) {
    let i = 0;
    element.innerHTML = ""; // 기존 텍스트 초기화

    function type() {
        // HTML 태그 감지 및 처리
        if (i < htmlText.length) {
            if (htmlText.charAt(i) === "<") {
                // 태그가 끝날 때까지 추가
                const tagEnd = htmlText.indexOf(">", i);
                if (tagEnd !== -1) {
                    element.innerHTML += htmlText.substring(i, tagEnd + 1);
                    i = tagEnd + 1;
                }
            } else {
                element.innerHTML += htmlText.charAt(i);
                i++;
            }
            setTimeout(type, speed);
        }
    }
    type();
}

// 초기 답변 생성 함수
function initialBotMessage() {
    const chatDisplay = document.getElementById("chat-display");

    const botMessage = document.createElement("p");
    botMessage.className = "bot";
    chatDisplay.appendChild(botMessage);

    const botText = "안녕하세요! 법률 상담과 관련된 질문을 입력해 주세요.";
    typeWriterEffect(botText, botMessage, 50);
}

// 페이지 로드 시 초기 안내 메시지 출력
document.addEventListener("DOMContentLoaded", function () {
    setTimeout(initialBotMessage, 2000);
});

// 메시지 전송 함수
function sendMessage() {
    const userInputValue = document.getElementById('user-input').value;
    if (userInputValue.trim() === "") return;

    addMessage("user", userInputValue);
    displayAllLawInfo();

    document.getElementById('user-input').value = '';
}

// 메시지 추가 함수
function addMessage(sender, message) {
    const chatDisplay = document.getElementById('chat-display');
    const messageElement = document.createElement('p');
    messageElement.classList.add(sender === "user" ? "user" : "bot");

    if (sender === "user") {
        const spanElement = document.createElement('span');
        spanElement.innerText = message;
        messageElement.appendChild(spanElement);
    } else {
        typeWriterEffect(message, messageElement); // 봇 메시지에 HTML 타자기 효과 적용
    }

    chatDisplay.appendChild(messageElement);
    chatDisplay.scrollTop = chatDisplay.scrollHeight;
}

// 모든 법률 정보를 출력하는 함수
function displayAllLawInfo() {
    const allLawInfo = `
딥페이크 관련 법적 조치 및 주요 법률<br><br>

1. 정보통신망 이용촉진 및 정보보호 등에 관한 법률 (정보통신망법)<br>
- 개요: 한국에서는 정보통신망법을 통해 딥페이크 영상을 허위로 생성하거나, 명예를 훼손하고 개인 정보를 침해하는 행위를 규제합니다. 특히, 딥페이크 영상으로 인해 발생하는 사이버 명예훼손 및 성적 수치심 유발 콘텐츠를 처벌합니다.
<a href="https://www.google.com/search?q=%EC%A0%95%EB%B3%B4%ED%86%B5%EC%8B%A0%EB%A7%9D+%EC%9D%B4%EC%9A%A9%EC%B4%89%EC%A7%84+%EB%B0%8F+%EC%A0%95%EB%B3%B4%EB%B3%B4%ED%98%B8+%EB%93%B1%EC%97%90+%EA%B4%80%ED%95%9C+%EB%B2%95%EB%A5%A0&oq=%EC%A0%95%EB%B3%B4%ED%86%B5%EC%8B%A0%EB%A7%9D+%EC%9D%B4%EC%9A%A9%EC%B4%89%EC%A7%84+%EB%B0%8F+%EC%A0%95%EB%B3%B4%EB%B3%B4%ED%98%B8+%EB%93%B1%EC%97%90+%EA%B4%80%ED%95%9C+%EB%B2%95%EB%A5%A0&gs_lcrp=EgZjaHJvbWUyCQgAEEUYORiABDIHCAEQABiABDIHCAIQABiABDIHCAMQABiABDIHCAQQABiABDIHCAUQABiABDIHCAYQABiABDIHCAcQABiABDIHCAgQABiABDIHCAkQABiABNIBBzE2OWowajeoAgCwAgA&sourceid=chrome&ie=UTF-8" target="_blank">정보통신망법 링크</a><br><br>

2. 성폭력범죄의 처벌 등에 관한 특례법 (성폭력처벌법)<br>
- 개요: 딥페이크 기술을 악용해 성적 이미지를 비동의로 합성하거나 유포하는 행위는 성폭력처벌법에 따라 처벌됩니다. 동의 없이 딥페이크 합성물을 유포하는 경우 최대 징역형까지 선고될 수 있습니다.<br>
<a href="https://www.google.com/search?q=%EC%84%B1%ED%8F%AD%EB%A0%A5%EB%B2%94%EC%A3%84%EC%9D%98+%EC%B2%98%EB%B2%8C+%EB%93%B1%EC%97%90+%EA%B4%80%ED%95%9C+%ED%8A%B9%EB%A1%80%EB%B2%95+(%EC%84%B1%ED%8F%AD%EB%A0%A5%EC%B2%98%EB%B2%8C%EB%B2%95)&oq=%EC%84%B1%ED%8F%AD%EB%A0%A5%EB%B2%94%EC%A3%84%EC%9D%98+%EC%B2%98%EB%B2%8C+%EB%93%B1%EC%97%90+%EA%B4%80%ED%95%9C+%ED%8A%B9%EB%A1%80%EB%B2%95+(%EC%84%B1%ED%8F%AD%EB%A0%A5%EC%B2%98%EB%B2%8C%EB%B2%95)&gs_lcrp=EgZjaHJvbWUyBggAEEUYOTIKCAEQABiABBiiBDIKCAIQABiABBiiBDIKCAMQABiABBiiBNIBBzE4M2owajeoAgiwAgE&sourceid=chrome&ie=UTF-8" target="_blank">성폭력처벌법 링크</a><br><br>

3. 저작권법<br>
- 개요: 원본 영상이나 이미지에 대한 저작권을 침해하는 딥페이크 콘텐츠는 저작권법 위반에 해당할 수 있습니다. 특히, 저작물의 무단 변형 및 복제, 배포 행위를 규제합니다.<br>
<a href="https://www.google.com/search?q=%EC%A0%80%EC%9E%91%EA%B6%8C%EB%B2%95&sca_esv=48f637e9bc078c4b&sxsrf=ADLYWIKlSIeid-hY8L5jT5sXFEPAhZhKhw%3A1730716380438&ei=3KIoZ-WrGrvT1e8P9e27-A8&ved=0ahUKEwjlicD4vMKJAxW7afUHHfX2Dv8Q4dUDCA8&uact=5&oq=%EC%A0%80%EC%9E%91%EA%B6%8C%EB%B2%95&gs_lp=Egxnd3Mtd2l6LXNlcnAiDOyggOyekeq2jOuylTILEAAYgAQYsQMYgwEyChAAGIAEGEMYigUyChAAGIAEGBQYhwIyChAAGIAEGEMYigUyCxAAGIAEGLEDGIMBMgoQABiABBhDGIoFMgUQABiABDIFEAAYgAQyBRAAGIAEMgUQABiABEjtBFDeA1jeA3ABeACQAQCYAYYBoAHzAaoBAzAuMrgBA8gBAPgBAvgBAZgCAqACkAHCAgoQABiwAxjWBBhHmAMAiAYBkAYKkgcDMS4xoAfADg&sclient=gws-wiz-serp" target="_blank">저작권법 링크</a><br><br>

4. 개인정보 보호법<br>
- 개요: 개인정보 보호법은 개인의 초상권 및 사생활 보호를 위해 딥페이크 영상에 포함된 개인 정보의 오남용을 규제합니다. 특정인의 얼굴을 사용한 합성 이미지와 영상을 생성하고 배포하는 행위는 개인정보 보호법 위반으로 처벌될 수 있습니다.<br>
<a href="https://www.google.com/search?q=%EA%B0%9C%EC%9D%B8%EC%A0%95%EB%B3%B4+%EB%B3%B4%ED%98%B8%EB%B2%95&sca_esv=48f637e9bc078c4b&sxsrf=ADLYWIJlLbE0Dw8DZoP8nXcVWOhhOQ8o8A%3A1730716402884&ei=8qIoZ7LXNZeavr0P7--F4QE&ved=0ahUKEwiymJqDvcKJAxUXja8BHe93IRwQ4dUDCA8&uact=5&oq=%EA%B0%9C%EC%9D%B8%EC%A0%95%EB%B3%B4+%EB%B3%B4%ED%98%B8%EB%B2%95&gs_lp=Egxnd3Mtd2l6LXNlcnAiFuqwnOyduOygleuztCDrs7TtmLjrspUyBRAAGIAEMgUQABiABDIFEAAYgAQyBRAAGIAEMgUQABiABDIKEAAYgAQYFBiHAjIFEAAYgAQyBRAAGIAEMgUQABiABDIFEAAYgARI0gRQkgNYkgNwAXgBkAEAmAFzoAFzqgEDMC4xuAEDyAEA-AEC-AEBmAICoAJ_wgIKEAAYsAMY1gQYR5gDAIgGAZAGCpIHAzEuMaAHyQg&sclient=gws-wiz-serp" target="_blank">개인정보 보호법 링크</a><br><br>

5. 형법 (명예훼손죄 및 모욕죄)<br>
- 개요: 딥페이크 영상이 특정인을 대상으로 허위 사실을 포함해 명예를 훼손하거나, 모욕하는 경우 형법상 명예훼손죄와 모욕죄가 적용될 수 있습니다.<br>
<a href="https://www.google.com/search?q=%ED%98%95%EB%B2%95+%28%EB%AA%85%EC%98%88%ED%9B%BC%EC%86%90%EC%A3%84%29&sca_esv=48f637e9bc078c4b&sxsrf=ADLYWIKBB6nBQhPVfu8zv6Ba2cUmWV231g%3A1730716436944&ei=FKMoZ4mnOdei1e8PhbbreA&ved=0ahUKEwiJgbmTvcKJAxVXUfUHHQXbGg8Q4dUDCA8&uact=5&oq=%ED%98%95%EB%B2%95+%28%EB%AA%85%EC%98%88%ED%9B%BC%EC%86%90%EC%A3%84%29&gs_lp=Egxnd3Mtd2l6LXNlcnAiGO2YleuylSAo66qF7JiI7Zu87IaQ7KOEKTIEEAAYHjIGEAAYCBgeMggQABgIGB4YDzIKEAAYCBgKGB4YDzIIEAAYgAQYogQyCBAAGIAEGKIEMggQABiABBiiBDIIEAAYCBgeGA9IywdQ1gJY5QVwAXgBkAEAmAF2oAHlAaoBAzAuMrgBA8gBAPgBAZgCA6AC7gHCAgcQIxiwAxgnwgIKEAAYsAMY1gQYR8ICBBAjGCeYAwCIBgGQBgqSBwMxLjKgB-EM&sclient=gws-wiz-serp" target="_blank">형법 링크</a><br><br>
`;

    addMessage("bot", allLawInfo);
}
