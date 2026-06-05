# Sprawozdanie Zadanie 2 - GitHub Actions

## Konfiguracja łańcucha (Pipeline)
Został opracowany łańcuch CI/CD, który realizuje następujące kroki:
1. **Budowa obrazu lokalnego:** Najpierw budowana jest wersja jednostkowa w celu przeprowadzenia audytu bezpieczeństwa.
2. **Skanowanie CVE (Trivy):** Wykorzystano skaner Trivy, który przerywa działanie łańcucha (zwracając `exit-code: 1`), jeśli wykryje podatności na poziomie `CRITICAL` lub `HIGH`.
3. **Budowa Multi-arch i Push:** Po pomyślnym przejściu skanowania, obraz jest budowany przy pomocy BuildKit dla architektur `linux/amd64` oraz `linux/arm64` i wysyłany do GitHub Container Registry (ghcr.io).
4. **Zarządzanie Cache:** Zewnętrzne dane cache są pobierane i wysyłane (w trybie `mode=max`) do dedykowanego publicznego repozytorium `zad2` na platformie Docker Hub.

## Strategia tagowania (Z uzasadnieniem)

**1. Tagowanie obrazu wynikowego (GHCR):**
Zastosowano podwójne tagowanie obrazów trafiających do głównego rejestru pod nazwą aplikacji:
* **`latest`**: Używany dla wygody, by wskazywać na najbardziej aktualną, stabilną wersję aplikacji (np. `zad2:latest`).
* **`${{ github.sha }}` (Git Commit SHA):** Każdy obraz otrzymuje również unikalny identyfikator hashujący konkretnego commita z repozytorium Git. 
*Uzasadnienie:* Zgodnie z dobrymi praktykami GitOps, tag `latest` jest tzw. tagiem mutowalnym (zmiennym) i nie nadaje się do bezpiecznych wdrożeń produkcyjnych, ponieważ nie gwarantuje przewidywalności (nie wiadomo precyzyjnie, z jakiego kodu powstał obraz). Tagowanie wartością `SHA` zapewnia absolutną identyfikowalność – pozwala w ułamku sekundy powiązać działający kontener z dokładną wersją kodu źródłowego w repozytorium.

**2. Tagowanie obrazu Cache (Docker Hub):**
* **`zad2:cache`**: Wszystkie metadane z procesu budowy wysyłane są pod jeden, stały tag.
*Uzasadnienie:* Tryb cache `max` eksportuje wszystkie warstwy dla wszystkich etapów budowy. Użycie stałego tagu `cache` sprawia, że przy kolejnych uruchomieniach środowisko CI/CD po prostu nadpisuje starą zawartość pamięci podręcznej. Zapobiega to drastycznemu rozrostowi zajmowanego miejsca w rejestrze Docker Hub, zachowując przy tym maksymalną wydajność budowania.