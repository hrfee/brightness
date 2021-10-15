#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#define PATH "BACKLIGHT_PATH"
#define MIN 0
#define MAX 7500
#define PCTMIN 2
#define PCTMAX 100
#define STEP MAX / 100
#define STEPS 20

float LOGMIN;
float LOGMAX;

int GetBrightness() {
  FILE *f = fopen(PATH, "r");
  if (f == NULL)
    return -1;
  int n;
  fscanf(f, "%i", &n);
  fclose(f);
  return n;
}

float GetBrightnessPct() {
  int b = GetBrightness();
  if (b == -1)
    return -1;
  return ((float)b / MAX) * 100;
}

int SetBrightness(int val) {
  if (val > MAX)
    return -1;
  FILE *f = fopen(PATH, "w");
  int e = fprintf(f, "%d", val);
  fclose(f);
  return e;
}

int SetBrightnessPct(float pct) { return SetBrightness((pct / 100) * MAX); }

int SetBrightnessSmooth(int val) {
  int current = GetBrightness();
  struct timespec wait;
  wait.tv_nsec = 1e7;
  if (current < 0)
    return current;
  if (val > current) {
    for (; current <= val; current += STEP) {
      SetBrightness(current);
      nanosleep(&wait, &wait);
    }
  } else {
    for (; current >= val; current -= STEP) {
      SetBrightness(current);
      nanosleep(&wait, &wait);
    }
  }
  return 0;
}

float BrightnessPctToStep(float pct) {
  return round(log10(pct) / (LOGMAX - LOGMIN) * STEPS);
}

float StepToBrightnessPct(float step) {
  float x = step / STEPS * (LOGMAX - LOGMIN);
  return round(fmax(fmin(pow(10, x), PCTMAX), PCTMIN));
}

void Usage(char *argv[]) {
  printf("usage\n\t%s -U / -B / -D\n", argv[0]);
  exit(1);
}

int main(int argc, char *argv[]) {
  LOGMIN = log10(PCTMIN);
  LOGMAX = log10(PCTMAX);

  if (argc < 2 || !(!strcmp(argv[1], "-U") || !strcmp(argv[1], "-D") ||
                    !strcmp(argv[1], "-B")))
    Usage(argv);

  float currentBrightness = GetBrightnessPct();
  if (currentBrightness < 0)
    printf("Failed to read file \"%s\": %i\n", PATH, (int)currentBrightness);

  if (currentBrightness == 0) {
    SetBrightnessPct(PCTMIN + 1);
    currentBrightness = GetBrightnessPct();
    if (currentBrightness < 0)
      printf("Failed to read file \"%s\": %i\n", PATH, (int)currentBrightness);
  }

  float currentStep = BrightnessPctToStep(currentBrightness);

  float newStep = currentStep;

  if (strcmp(argv[1], "-D") == 0)
    newStep -= 2;
  else if (strcmp(argv[1], "-U") == 0)
    newStep += 2;
  else if (strcmp(argv[1], "-B") == 0)
    newStep += 5;

  float newBrightness = StepToBrightnessPct(newStep);

  printf("Current backlight: %.1f\nChanging to %.1f\n", currentBrightness,
         newBrightness);

  if (newBrightness == 0)
    exit(1);

  if (strcmp(argv[1], "-U") && currentBrightness == 99)
    newBrightness = 100;

  int e = SetBrightnessSmooth((newBrightness / 100) * MAX);
  if (e < 0) {
    printf("Failed to write brightness: %i\n", e);
    exit(1);
  }
}
