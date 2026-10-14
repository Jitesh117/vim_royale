![Demo](./assets/vim_royale_hero.png)
<div align="center">
    <a href="https://peerlist.io/jitesh117/project/vim-royale" target="_blank" rel="noreferrer">
        <img
            src="https://peerlist.io/api/v1/projects/embed/PRJHKKD9N7D9OAEPGIDJE68EPMLNAN?showUpvote=true&theme=light"
            alt="Vim Royale"
            style="width: auto; height: 72px;"
        />
    </a>
</div>


## About

Vim Royale is a pvp vim based multiplayer game. Users can play 1v1 matches against each other and compete for a spot on the leaderboard. The game uses AI to generate random code snippets for players to edit and compete on. Users can sign in with Google to save their progress and compete on the leaderboard. 

## Features
- 1v1 matches with real-time updates
- Rating system based on Elo
- Spectator mode for matches
- Create and join custom lobbies with your friends
- Match replay to analyse your moves
- Built-in VimTutor
- Leaderboard to compete for a spot
- Code snippets for players to edit and compete on

## Tech Stack

- Frontend: React, TypeScript, Vite, CodeMirror
- Backend: Go, Gin, WebSockets, GORM
- Database: PostgreSQL
- Auth: Google/GitHub OAuth

## Want to Contribute?

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md)

## Security

Refer to [SECURITY.md](./SECURITY.md)

## Running Locally

```bash
cd backend && docker compose -f docker-compose.dev.yml up --build   
```

```bash
cd frontend && npm install && npm run dev
```

## Support

If you wish to support further development of Vim Royale, you can donate via [GitHub Sponsors](https://github.com/sponsors/Jitesh117), [BuyMeACoffee](https://www.buymeacoffee.com/jitesh117) or [PayPal](https://www.paypal.com/paypalme/jitesh1177)

## Credits

- [Primeagen](https://github.com/ThePrimeagen) for the idea of the game and the name via his streams.
