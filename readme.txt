Nakonec byl tento projekt dělán jen jedním studentem, takže se nic nerozdělovalo
Proto jsem také nepsal hlavičky souborů

Celý projekt je také na githubu https://github.com/Siigull/ITU

Backend 
 - Je potřeba mít nainstalovaný jazyk go. Projekt byl vyvíjen na verzi 1.23 za jiné neručím
 - Složka Backend obsahuje implementaci backend serveru uvnitř main.go
 - Také obsahuje příklady volání některých api endpointů
 - Pro spuštění je potřeba přejít do backend složky a použít go run .
   - Poté je jeho endpoint na localhost:8080

Frontend
 - Je potřeba nainstalovat node a potom next přes npm (npm install next)
 - Kód kterej jsem vytvářel je jen uvnitř src
   -> Stránka má tři sekce
      - uvodní (volba přezdívky) page.js v kořenové
      - menu Zobrazuje založené hry a všechny uživatele, je v menu složce
      - hra Zobrazuje vybranou hru, hlavní sekce, je v game složce
   
   -> Některé komponenty jsou mimo kód sekce, jsou uloženy v kořenové, jsx koncovka
   -> Funkce co se používaj ve více souborech jsou v helper.js
      - funkce pro reaktivní proměnnou na rozměry okna je vzána ze stack overflow [https://stackoverflow.com/questions/73070114/i-want-to-change-style-according-to-window-width-using-states]

 - Frontend se nejjednodušeji spustí přes npm run dev
   - Poté je na localhost:3000