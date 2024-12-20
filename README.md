# BubbleBubble (🥇 과기부 | 민관협력 플랫폼 특별상)

![11](https://github.com/user-attachments/assets/fd6e8ea5-2c2d-4f65-8632-1c32e83c5291)


1. 프로젝트 개요<br>
버블버블은 네이버 클라우드 챗봇과 AI를 활용하여 신종 범죄 동향과 사용자 기반 맞춤형 정보를 제공하는 서비스입니다. <br>
특히 실시간 키워드 분석을 통해 관련된 주요 키워드를 추천하고, 24시간 자동화된 챗봇 서비스, 맞춤형 학습 기능, <br>
통합 생성형 지식 서비스(GPT)를 제공합니다.<br>

2. 주요 기능<br>
신종 범죄 사용자 정보 제공: 사용자 정보를 바탕으로 최신 범죄 동향과 관련 정보를 제공<br><br>
키워드분석 및 추천: 실시간 키워드 분석을 통해 사용자에게 맞춤형 키워드를 추천하여 최신 정보를 빠르게 전달<br>
복지향상 챗봇(버블버블 챗봇): 사용자 정보 기반으로 24시간 서비스 자동화를 지원하여 사용자의 편의성<br>
맞춤형 학습 서비스(버블버블 퀴즈): 사용자 정보에 맞춘 개인화된 학습 퀴즈 기능을 제공<br>
통합생성 지식서비스(버블버블 GPT): 사용자 정보 기반으로 지식과 정보를 통합하여 제공하는 AI 지식 서비스<br>

3. 폴더구조<br> 
```
├── public/                 # 정적 파일 및 아이콘
│   ├── index.html          # 메인 HTML 파일
│   └── manifest.json       # PWA 설정 파일
├── src/
│   ├── assets/             # 이미지 및 폰트 파일
│   ├── components/         # UI 컴포넌트 파일
│   │   ├── Chatbot.js      # 법률 챗봇 UI 및 로직
│   │   ├── GameService.js  # 법률 퀴즈 및 게임 기능
│   │   ├── GPTService.js   # 통합 GPT 지식 서비스
│   │   └── Keyword.js      # 키워드 추천 기능
│   ├── services/           # API 호출 및 데이터 로직
│   │   ├── naverCloudAPI.js # 네이버 클라우드 API 연동
│   ├── styles/             # CSS 파일 및 스타일
│   ├── App.js              # 메인 앱 컴포넌트
│   └── index.js            # 진입점 파일
└── README.md
```

Copyright by @bubblebubble-labs
