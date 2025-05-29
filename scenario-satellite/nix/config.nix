{ lib, ... }:
{
  options =
    with lib;
    with types;
    {
      # Public and private ssh key of the web server
      web_public = mkOption { type = str; };
      web_private = mkOption { type = str; };

      sh_messages = mkOption { type = str; };

      # IP of DNS autoritary server for c2
      autoritary = mkOption { type = str; };
      # IP of DNS resolver
      resolver = mkOption { type = str; };
      # IP of the satellite
      satellite = mkOption { type = str; };
      # IP of the web server used in ip table
      web = mkOption { type = str; };
    };

  config = {
    web_public = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKe9aUD5W9n5Uu7izwbQAkpWPKyouk4E3XhM/D8vOAh2 admin@nixos";
    web_private = ''
      -----BEGIN OPENSSH PRIVATE KEY-----
      b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
      QyNTUxOQAAACCnvWlA+VvZ+VLu4s8G0AJKVjysqLpOBN14TPw/LzgIdgAAAJBQFJm1UBSZ
      tQAAAAtzc2gtZWQyNTUxOQAAACCnvWlA+VvZ+VLu4s8G0AJKVjysqLpOBN14TPw/LzgIdg
      AAAEDWJZNplqQhZCADKIzCZfqEb5tfg9Xw3wtxx6MQZWdFG6e9aUD5W9n5Uu7izwbQAkpW
      PKyouk4E3XhM/D8vOAh2AAAAC2FkbWluQG5peG9zAQI=
      -----END OPENSSH PRIVATE KEY-----
    '';

    sh_messages = # bash
      ''
        # Array of messages to be sent
        messages=(
            "Hello, world!"
            "This is a test message."
            "Sending data via HTTP POST."
            "Random message content."
            "Another message in the loop."
        )

        # URL to send the POST request to
        url="http://satellite.app/"

        # Function to log messages
        log() {
            echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
        }

        # Loop indefinitely
        while true; do
            # Select a random message from the array
            message="''${messages[$RANDOM % ''${#messages[@]}]}"

            # Send the POST request with the message
            if curl -X POST \
                    -H "Content-Type: application/json" \
                    -d "{\"message\":\"$message\"}" \
                    --silent \
                    --output /dev/null \
                    --connect-timeout 2 \
                    --max-time 2 \
                    "$url"; then
                log "Successfully sent message: $message"
            else
                log "Failed to send message: $message"
            fi

            # Wait for a random number of seconds between 1 and 5
            sleep $((RANDOM % 5 + 1))
        done
      '';

    autoritary = "10.0.2.2@8054";
    resolver = "10.0.2.2:8053";
    web = "10.0.2.2:8080";
    # To get difference between IP destination for a same machine, use interface
    satellite = "10.0.2.2";
  };
}
