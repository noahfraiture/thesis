#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#define SERVER_IP "127.0.0.1"
#define SERVER_PORT 1881

struct robot_control {
  char id[16];
  int authorized;
};

int send_message(const char *message) {
  printf("Send message : %s\n", message);
  int sock = socket(AF_INET, SOCK_DGRAM, 0);
  if (sock < 0) {
    perror("Socket creation failed");
    return 1;
  }

  struct sockaddr_in server_addr;
  memset(&server_addr, 0, sizeof(server_addr));
  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(SERVER_PORT);
  if (inet_pton(AF_INET, SERVER_IP, &server_addr.sin_addr) <= 0) {
    perror("Invalid address");
    close(sock);
    return 1;
  }
  if (sendto(sock, message, strlen(message), 0, (struct sockaddr *)&server_addr,
             sizeof(server_addr)) < 0) {
    perror("Send failed");
    return 1;
  }
  close(sock);
  return 0;
}

int execute_command(char *command) {
  if (strcmp(command, "quit") == 0) {
    return 1;
  }

  if (strcmp(command, "forward") == 0) {
    const char *message = "forward";

    if (send_message("robot")) {
      perror("Send failed");
      return 1;
    }
  } else {
    printf("Command '%s' accepted but not sent (only 'forward' sends UDP).\n",
           command);
  }
  return 0;
}

int robot_control() {
  struct robot_control rc;
  rc.authorized = 0;

  printf("Available robots:\n");
  printf("1. RobotA\n");
  printf("2. RobotB\n");
  printf("3. RobotC\n");

  printf("Enter the ID of the robot you want to control: ");
  gets(rc.id);

  if (rc.authorized) {
    printf("Authorization granted. You can now send commands.\n");

    char command[10];
    while (1) {
      printf("Enter command (left, forward, back, right, or quit): ");
      if (scanf("%9s", command) < 0) {
        return 1;
      }
      if (execute_command(command)) {
        break;
      }
    }

  } else {
    printf("Authorization denied. You cannot control this robot.\n");
  }
  return 0;
}

int lock_control() {
  char input[50];
  char password[] = {0x31, 0x27, 0x21, 0x30,
                     0x27, 0x36, 0x00}; // 'secret' obfuscated
  char key = 0x42;
  char deobfuscated[7];

  // deobfuscated password
  strcpy(deobfuscated, password);
  for (int i = 0; i < strlen(deobfuscated); i++) {
    deobfuscated[i] ^= key;
  }

  // compare password
  printf("Enter password: ");
  if (fgets(input, sizeof(input), stdin) == NULL) {
    return 1;
  }
  input[strcspn(input, "\n")] = 0;
  if (strcmp(input, deobfuscated) == 0) {
    printf("Access granted, switch lock!\n");
    return send_message("door");
  } else {
    printf("Access denied!\n");
  }
  return 0;
}

int main(int argc, char *argv[]) {
  if (strstr(argv[0], "robot") != NULL) {
    return robot_control();
  } else if (strstr(argv[0], "lock") != NULL) {
    return lock_control();
  } else {
    printf("%s", argv[0]);
  }
  return 0;
}
