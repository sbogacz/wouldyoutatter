#!/usr/bin/env python3

""" Simple test of using the API from a different language
    (i.e. not depending on the predefined types).

"""

import requests
import random
import json

ENDPOINT = 'https://rl8b1iyh90.execute-api.us-west-2.amazonaws.com/production'

def main():

    jar = requests.cookies.RequestsCookieJar()

    for _ in range(100):

        # Get a random matchup. Uses the cookie jar to keep track of matchups
        # that have been seen. A browser would handle this for us.
        r = requests.get(ENDPOINT + '/matchups/random', cookies=jar)
        matchup = r.json()

        # Biased auto-voter
        winner = None
        for contender in ('contender_1', 'contender_2'):
            if matchup[contender]['name'] == 'steve':
                winner = contender
                break

        if not winner:
            # Randomly choose a contender
            if random.random() > 0.5:
                winner = 'contender_1'
            else:
                winner = 'contender_2'

        vote = {'Winner': matchup[winner]['name']}
        print("Vote Payload: {}".format(vote))

        r = requests.post(ENDPOINT + matchup['vote_url'], json=vote, cookies=jar)

    print("\nCurrent leaderboard:")
    r = requests.get(ENDPOINT + '/leaderboard?limit=10')
    leaders = r.json()
    for leader in leaders:
        print("{}: Score: {} Wins: {} Losses: {}".format(leader['name'], leader['score'], leader['wins'], leader['losses']))


if __name__ == '__main__':
    main()
