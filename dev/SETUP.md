# Dev Env 

## Initial Local env setup

0. Duplicate `docker-compose-dev-env.example.yml` and rename it to `docker-compose-dev-env.yml`

1. Run `./dev.sh` script found in this directory. Wait until everything is up and running 

2. Create a new Organzation:
```bash
$ docker exec -u zulip go-zulip-zulip-1 /home/zulip/deployments/current/manage.py generate_realm_creation_link

Please visit the following secure single-use link to register your new Zulip organization:

https://localhost/new/jj46mgxl7cbxazmoi7evtxjx
```

3. Get get your API Key
    1. Go to http://localhost (accept to `Proceed to localhost (unsafe)`)
    2. Profile > Settings > Account & Privacy > `API Key`

4. Set up the env vars in `docker-compose-dev-env.yml` adding your registered email and the `API Key`

5. Give permissions to allow to the owner to create Users via API
    ```shell
    $ docker exec -u zulip go-zulip-zulip-1 /home/zulip/deployments/current/manage.py change_user_role your_user_email@test.com can_create_users -r 2

    User can already create users for this realm.
    ```

6. `Ctrl+C` to stop all

## Developing

1. Run `./dev.sh` script
2. Go to `https://localhost` to access to Zulip
3. Develop, run integration tests, check the results
4. A filewatcher (`air`) will run linters and tests on every `.go` file change running linters, test and a small example application.
