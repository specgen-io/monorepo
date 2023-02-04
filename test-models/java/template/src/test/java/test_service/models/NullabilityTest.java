package test_service.models;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

public class NullabilityTest {
  @Test
  public void oneOfItemNotNull() {
    assertThrows(IllegalArgumentException.class, () -> new OrderEventWrapper.Canceled(null));
  }
}
