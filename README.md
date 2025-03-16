# Spotify API

## À propos du projet
Spotify API est une application web qui utilise l'API Spotify pour permettre aux utilisateurs de découvrir de nouvelles musiques, de rechercher des artistes et des albums, et de gérer leurs favoris. Notre application propose diverses fonctionnalités telles que la découverte aléatoire de titres, la gestion d'un compte utilisateur, et la visualisation détaillée des informations sur les titres.

## Fonctionnalités
- Découverte aléatoire de titres musicaux
- Système d'authentification (inscription et connexion)
- Gestion de favoris
- Recherche de titres, artistes et albums
- Affichage détaillé des informations sur les titres

## Installation

### Prérequis
- Go (version 1.15 ou supérieure)
- Accès à l'API Spotify (clé d'API)

### Étapes d'installation
1. Clonez le dépôt :
```bash git
clone https://github.com/jfrbast/Spotify.git
```
2. Configurez vos identifiants Spotify dans un fichier de configuration (non inclus dans ce dépôt pour des raisons de sécurité).

3. Lancez l'application :
```bash
go run main.go
```
4. Accédez à l'application via votre navigateur à l'adresse `http://localhost:8080`

## Routes implémentées

| Route | Méthode | Description |
|-------|---------|-------------|
| `/` | GET | Page d'accueil |
| `/random` | GET | Affichage d'un titre aléatoire |
| `/random/treatment` | POST | Traitement après sélection d'un titre aléatoire |
| `/random/remove` | POST | Suppression d'un titre des favoris depuis la page random |
| `/header` | GET | Gestion de l'en-tête du site |
| `/detail` | GET | Affichage détaillé d'un titre |
| `/recherche` | GET | Recherche de titres, artistes ou albums |
| `/compte` | GET | Gestion du compte utilisateur |
| `/compte/removefavorite` | POST | Suppression d'un titre des favoris depuis la page compte |
| `/compte/deconnexion` | GET | Déconnexion de l'utilisateur |
| `/inscription` | GET | Page d'inscription |
| `/inscription/treatment` | POST | Traitement du formulaire d'inscription |
| `/connexion` | GET | Page de connexion |
| `/connexion/treatment` | POST | Traitement du formulaire de connexion |

## API Spotify

Nous utilisons l'API Spotify pour accéder aux données suivantes :
- Informations sur les titres (nom, artiste, album, durée, popularité)
- Informations sur les artistes (nom, genres, popularité)
- Informations sur les albums (nom, date de sortie, nombre de titres)
- Recherche de titres, artistes et albums

### Endpoints utilisés
- `GET /v1/tracks/{id}` - Récupération des informations d'un titre
- `GET /v1/search` - Recherche de titres, artistes ou albums
- `GET /v1/artists/{id}` - Récupération des informations d'un artiste
- `GET /v1/albums/{id}` - Récupération des informations d'un album
- `GET /v1/recommendations` - Recommandations basées sur des seeds

## Base de données
Notre application utilise une base de données pour stocker :
- Les informations des utilisateurs (nom, email, mot de passe haché)
- Les titres favoris des utilisateurs
- Les historiques de recherche et d'écoute

## Synthèse du déroulement du projet

### Décomposition du projet
Nous avons décomposé le projet en plusieurs phases clés :
1. **Phase de conception** : Définition des fonctionnalités et de l'architecture
2. **Phase de développement backend** : Mise en place de l'API, connexion à Spotify, gestion utilisateurs
3. **Phase de développement frontend** : Création des templates et de l'interface utilisateur
4. **Phase de tests** : Vérification du bon fonctionnement de l'application
5. **Phase de déploiement** : Préparation et mise en ligne de l'application


 API Spotify et intégration
 Système d'authentification et gestion des utilisateurs
Interface utilisateur et templates
 Tests et documentation

Nous avons utilisé une approche Agile avec des sprints de deux semaines et des réunions hebdomadaires pour suivre l'avancement du projet.

### Gestion du temps
Pour gérer efficacement notre temps, nous avons :
- Défini des priorités en nous concentrant d'abord sur les fonctionnalités essentielles
- Établi un calendrier avec des jalons clairs
- Utilisé un outil de suivi de projet (Trello) pour visualiser l'avancement
- Prévu des marges pour les imprévus et les difficultés techniques

### Documentation
Notre stratégie de documentation a consisté à :
- Consulter régulièrement la documentation officielle de l'API Spotify
- Participer à des forums spécialisés (Stack Overflow, Reddit)
- Échanger avec d'autres développeurs ayant déjà travaillé avec l'API Spotify
- Créer une documentation interne pour faciliter la collaboration entre membres

## Conclusion
Ce projet nous a permis de mettre en pratique nos connaissances en développement web avec Go, d'apprendre à utiliser une API externe et de travailler en équipe sur un projet complet. Les défis principaux ont été la gestion de l'authentification Spotify et l'optimisation des requêtes API.